package controller

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/snebel29/kooper/log"
	"github.com/snebel29/kooper/monitoring/metrics"
	"github.com/snebel29/kooper/operator/controller/leaderelection"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kooper/operator/common"
)

// Span tag and log keys.
const (
	kubernetesObjectKeyKey         = "kubernetes.object.key"
	kubernetesObjectNSKey          = "kubernetes.object.namespace"
	kubernetesObjectNameKey        = "kubernetes.object.name"
	eventKey                       = "event"
	kooperControllerKey            = "kooper.controller"
	processedTimesKey              = "kubernetes.object.total_processed_times"
	retriesRemainingKey            = "kubernetes.object.retries_remaining"
	processingRetryKey             = "kubernetes.object.processing_retry"
	retriesExecutedKey             = "kubernetes.object.retries_consumed"
	controllerNameKey              = "controller.cfg.name"
	controllerResyncKey            = "controller.cfg.resync_interval"
	controllerMaxRetriesKey        = "controller.cfg.max_retries"
	controllerConcurrentWorkersKey = "controller.cfg.concurrent_workers"
	successKey                     = "success"
	messageKey                     = "message"
)

// generic controller is a controller that can be used to create different kind of controllers.
type generic struct {
	queue     workqueue.RateLimitingInterface // queue will have the jobs that the controller will get and send to handlers.
	informer  cache.SharedIndexInformer       // informer will notify be inform us about resource changes.
	handler   handler.Handler                 // handler is where the logic of resource processing.
	running   bool
	runningMu sync.Mutex
	cfg       Config
	tracer    opentracing.Tracer // use directly opentracing API because it's not an implementation.
	metrics   metrics.Recorder
	leRunner  leaderelection.Runner
	logger    log.Logger
	hasSynced bool
}

// NewSequential creates a new controller that will process the received events sequentially.
// This constructor is just a wrapper to help bootstrapping default sequential controller.
func NewSequential(resync time.Duration, handler handler.Handler, retriever retrieve.Retriever, metricRecorder metrics.Recorder, logger log.Logger) Controller {
	cfg := &Config{
		ConcurrentWorkers: 1,
		ResyncInterval:    resync,
	}
	return New(cfg, handler, retriever, nil, nil, metricRecorder, logger)
}

// NewConcurrent creates a new controller that will process the received events concurrently.
// This constructor is just a wrapper to help bootstrapping default concurrent controller.
func NewConcurrent(concurrentWorkers int, resync time.Duration, handler handler.Handler, retriever retrieve.Retriever, metricRecorder metrics.Recorder, logger log.Logger) (Controller, error) {
	if concurrentWorkers < 2 {
		return nil, fmt.Errorf("%d is not a valid concurrency workers ammount for a concurrent controller", concurrentWorkers)
	}

	cfg := &Config{
		ConcurrentWorkers: concurrentWorkers,
		ResyncInterval:    resync,
	}
	return New(cfg, handler, retriever, nil, nil, metricRecorder, logger), nil
}

// New creates a new controller that can be configured using the cfg parameter.
func New(cfg *Config, handler handler.Handler, retriever retrieve.Retriever, leaderElector leaderelection.Runner, tracer opentracing.Tracer, metricRecorder metrics.Recorder, logger log.Logger) Controller {
	// Sets the required default configuration.
	cfg.setDefaults()

	// Default logger.
	if logger == nil {
		logger = &log.Std{}
		logger.Warningf("no logger specified, fallback to default logger, to disable logging use dummy logger")
	}

	// Default metrics recorder.
	if metricRecorder == nil {
		metricRecorder = metrics.Dummy
		logger.Warningf("no metrics recorder specified, disabling metrics")
	}

	// Default tracer.
	if tracer == nil {
		tracer = &opentracing.NoopTracer{}
	}

	// If no name on controller do our best to infer a name based on the handler.
	if cfg.Name == "" {
		cfg.Name = reflect.TypeOf(handler).String()
		logger.Warningf("controller name not provided, it should have a name, fallback name to: %s", cfg.Name)
	}

	// Create the queue that will have our received job changes. It's rate limited so we don't have problems when
	// a job processing errors every time is processed in a loop.
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// store is the internal cache where objects will be store.
	store := cache.Indexers{}
	informer := cache.NewSharedIndexInformer(retriever.GetListerWatcher(), retriever.GetObject(), cfg.ResyncInterval, store)

	// Create our generic controller object.
	g := &generic{
		queue:     queue,
		informer:  informer,
		logger:    logger,
		hasSynced: false,
		metrics:   metricRecorder,
		tracer:    tracer,
		handler:   handler,
		leRunner:  leaderElector,
		cfg:       *cfg,
	}

	// Set up our informer event handler.
	// Objects are already in our local store. Add only keys/jobs on the queue so they can bre processed
	// afterwards.
	informer.AddEventHandlerWithResyncPeriod(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(&common.K8sEvent{Key: key, HasSynced: g.hasSynced, Kind: "Add"})
				metricRecorder.IncResourceEventQueued(cfg.Name, metrics.AddEvent)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(&common.K8sEvent{Key: key, HasSynced: g.hasSynced, Kind: "Update"})
				metricRecorder.IncResourceEventQueued(cfg.Name, metrics.AddEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(&common.K8sEvent{Key: key, HasSynced: g.hasSynced, Kind: "Delete"})
				metricRecorder.IncResourceEventQueued(cfg.Name, metrics.DeleteEvent)
			}
		},
	}, cfg.ResyncInterval)

	return g
}

func (g *generic) isRunning() bool {
	g.runningMu.Lock()
	defer g.runningMu.Unlock()
	return g.running
}

func (g *generic) setRunning(running bool) {
	g.runningMu.Lock()
	defer g.runningMu.Unlock()
	g.running = running
}

// Run will run the controller.
func (g *generic) Run(stopC <-chan struct{}) error {
	// Check if leader election is required.
	if g.leRunner != nil {
		return g.leRunner.Run(func() error {
			return g.run(stopC)
		})
	}

	return g.run(stopC)
}

// run is the real run of the controller.
func (g *generic) run(stopC <-chan struct{}) error {
	if g.isRunning() {
		return fmt.Errorf("controller already running")
	}

	g.logger.Infof("starting controller")
	// Set state of controller.
	g.setRunning(true)
	defer g.setRunning(false)

	// Shutdown when Run is stopped so we can process the last items and the queue doesn't
	// accept more jobs.
	defer g.queue.ShutDown()

	// Run the informer so it starts listening to resource events.
	go g.informer.Run(stopC)

	// Wait until our store, jobs... stuff is synced (first list on resource, resources on store and jobs on queue).
	if !cache.WaitForCacheSync(stopC, g.informer.HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}
	g.hasSynced = true

	// Start our resource processing worker, if finishes then restart the worker. The workers should
	// not end.
	for i := 0; i < g.cfg.ConcurrentWorkers; i++ {
		go func() {
			wait.Until(g.runWorker, time.Second, stopC)
		}()
	}

	// Until will be running our workers in a continuous way (and re run if they fail). But
	// when stop signal is received we must stop.
	<-stopC
	g.logger.Infof("stopping controller")

	return nil
}

// runWorker will start a processing loop on event queue.
func (g *generic) runWorker() {
	for {
		// Process newxt queue job, if needs to stop processing it will return true.
		if g.getAndProcessNextJob() {
			break
		}
	}
}

// getAndProcessNextJob job will process the next job of the queue job and returns if
// it needs to stop processing.
func (g *generic) getAndProcessNextJob() bool {
	// Get next job.
	nextJob, exit := g.queue.Get()
	if exit {
		return true
	}
	defer g.queue.Done(nextJob)
	evt := nextJob.(*common.K8sEvent)

	// Our root span will start here.
	span := g.tracer.StartSpan("processJob")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	g.setRootSpanInfo(evt.Key, span)

	// Process the job. If errors then enqueue again.
	if err := g.processJob(ctx, evt); err == nil {
		g.queue.Forget(evt)
		g.setForgetSpanInfo(evt.Key, span, err)
	} else if g.queue.NumRequeues(evt) < g.cfg.ProcessingJobRetries {
		// Job processing failed, requeue.
		g.logger.Warningf("error processing %s job (requeued): %v", evt.Key, err)
		g.queue.AddRateLimited(evt)
		g.metrics.IncResourceEventQueued(g.cfg.Name, metrics.RequeueEvent)
		g.setReenqueueSpanInfo(evt.Key, span, err)
	} else {
		g.logger.Errorf("Error processing %s: %v", evt.Key, err)
		g.queue.Forget(evt)
		g.setForgetSpanInfo(evt.Key, span, err)
	}

	return false
}

// processJob is where the real processing logic of the item is.
func (g *generic) processJob(ctx context.Context, evt *common.K8sEvent) error {
	// Get the object
	obj, exists, err := g.informer.GetIndexer().GetByKey(evt.Key)
	if err != nil {
		return err
	}

	// handle the object.
	if !exists { // Deleted resource from the cache.
		return g.handleDelete(ctx, evt)
	}

	evt.Object = obj.(runtime.Object)
	return g.handleAdd(ctx, evt)
}

func (g *generic) handleAdd(ctx context.Context, evt *common.K8sEvent) error {
	start := time.Now()
	g.metrics.IncResourceEventProcessed(g.cfg.Name, metrics.AddEvent)
	defer func() {
		g.metrics.ObserveDurationResourceEventProcessed(g.cfg.Name, metrics.AddEvent, start)
	}()

	// Create the span.
	pSpan := opentracing.SpanFromContext(ctx)
	span := g.tracer.StartSpan("handleAddObject", opentracing.ChildOf(pSpan.Context()))
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	// Set span data.
	ext.SpanKindConsumer.Set(span)
	span.SetTag(kubernetesObjectKeyKey, evt.Key)
	g.setCommonSpanInfo(span)
	span.LogKV(
		eventKey, "add",
		kubernetesObjectKeyKey, evt.Key,
	)

	// Handle the job.
	if err := g.handler.Add(ctx, evt); err != nil {
		ext.Error.Set(span, true) // Mark error as true.
		span.LogKV(
			eventKey, "error",
			messageKey, err,
		)

		g.metrics.IncResourceEventProcessedError(g.cfg.Name, metrics.AddEvent)
		return err
	}
	return nil
}

func (g *generic) handleDelete(ctx context.Context, evt *common.K8sEvent) error {
	start := time.Now()
	g.metrics.IncResourceEventProcessed(g.cfg.Name, metrics.DeleteEvent)
	defer func() {
		g.metrics.ObserveDurationResourceEventProcessed(g.cfg.Name, metrics.DeleteEvent, start)
	}()

	// Create the span.
	pSpan := opentracing.SpanFromContext(ctx)
	span := g.tracer.StartSpan("handleDeleteObject", opentracing.ChildOf(pSpan.Context()))
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	// Set span data.
	ext.SpanKindConsumer.Set(span)
	span.SetTag(kubernetesObjectKeyKey, evt.Key)
	g.setCommonSpanInfo(span)
	span.LogKV(
		eventKey, "delete",
		kubernetesObjectKeyKey, evt.Key,
	)

	// Handle the job.
	if err := g.handler.Delete(ctx, evt); err != nil {
		ext.Error.Set(span, true) // Mark error as true.
		span.LogKV(
			eventKey, "error",
			messageKey, err,
		)

		g.metrics.IncResourceEventProcessedError(g.cfg.Name, metrics.DeleteEvent)
		return err
	}
	return nil
}

func (g *generic) setCommonSpanInfo(span opentracing.Span) {
	ext.Component.Set(span, "kooper")
	span.SetTag(kooperControllerKey, g.cfg.Name)
	span.SetTag(controllerNameKey, g.cfg.Name)
	span.SetTag(controllerResyncKey, g.cfg.ResyncInterval)
	span.SetTag(controllerMaxRetriesKey, g.cfg.ProcessingJobRetries)
	span.SetTag(controllerConcurrentWorkersKey, g.cfg.ConcurrentWorkers)
}

func (g *generic) setRootSpanInfo(key string, span opentracing.Span) {
	numberRetries := g.queue.NumRequeues(key)

	// Try to set the namespace and resource name.
	if ns, name, err := cache.SplitMetaNamespaceKey(key); err == nil {
		span.SetTag(kubernetesObjectNSKey, ns)
		span.SetTag(kubernetesObjectNameKey, name)
	}

	g.setCommonSpanInfo(span)
	span.SetTag(kubernetesObjectKeyKey, key)
	span.SetTag(processedTimesKey, numberRetries+1)
	span.SetTag(processingRetryKey, numberRetries > 0)
	span.SetBaggageItem(kubernetesObjectKeyKey, key)
	ext.SpanKindConsumer.Set(span)
	span.LogKV(
		eventKey, "process_object",
		kubernetesObjectKeyKey, key,
	)
}

func (g *generic) setReenqueueSpanInfo(key string, span opentracing.Span, err error) {
	// Mark root span with error.
	ext.Error.Set(span, true)
	span.LogKV(
		eventKey, "error",
		messageKey, err,
	)

	rt := g.queue.NumRequeues(key)
	span.LogKV(
		eventKey, "reenqueued",
		retriesRemainingKey, g.cfg.ProcessingJobRetries-rt,
		retriesExecutedKey, rt,
		kubernetesObjectKeyKey, key,
	)
	span.LogKV(successKey, false)
}

func (g *generic) setForgetSpanInfo(key string, span opentracing.Span, err error) {
	success := true
	message := "object processed correctly"

	// Error data.
	if err != nil {
		// Mark root span with error.
		ext.Error.Set(span, true)
		span.LogKV(
			eventKey, "error",
			messageKey, err,
		)
		success = false
		message = "max number of retries reached after failing, forgetting object key"
	}

	span.LogKV(
		eventKey, "forget",
		messageKey, message,
		kubernetesObjectKeyKey, key,
	)
	span.LogKV(successKey, success)
}

