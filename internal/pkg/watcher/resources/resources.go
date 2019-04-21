package resources

import (
	"context"
	"errors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Event struct {
	key       string
	action    string
	namespace string
	kind      string
}

type ResourceWatcher interface {
	Run() error
	Shutdown()
	HasSynced() bool
}

type K8sResourceWatcher struct {
	ctx      context.Context
	kind     string
	informer cache.SharedInformer
	queue    workqueue.RateLimitingInterface
}

/*
func (r *K8sResourceWatcher) PopItem() {
	newEvent, quit := r.queue.Get()

	if quit {
		return false
	}
	defer r.queue.Done(newEvent)
	err := r.processItem(newEvent.(Event))
	if err == nil {
		// No error, reset the ratelimit counters
		r.queue.Forget(newEvent)
	} else if r.queue.NumRequeues(newEvent) < maxRetries {
		r.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		r.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		r.logger.Errorf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		r.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}
*/

func (r *K8sResourceWatcher) Run() error {
	stopCh := make(chan struct{})
	defer r.Shutdown(stopCh)

	go r.informer.Run(stopCh)
	if !cache.WaitForCacheSync(stopCh, r.HasSynced) {
		return errors.New("Timeout waiting for caches to sync")
	}
	return nil
}

func (r *K8sResourceWatcher) Shutdown(stopCh chan<- struct{}) {
	r.queue.ShutDown()
	close(stopCh)
}

func (r *K8sResourceWatcher) HasSynced() bool {
	return r.informer.HasSynced()
}
