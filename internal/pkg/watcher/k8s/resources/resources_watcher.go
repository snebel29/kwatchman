package resources

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/snebel29/kooper/monitoring/metrics"
	"github.com/snebel29/kooper/operator/controller"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

// K8sResourceWatcher represent the resourceWatcher
type K8sResourceWatcher struct {
	kind  string
	stopC chan struct{}
	ctrl  controller.Controller
}

// Run the resource watcher
func (r *K8sResourceWatcher) Run() error {
	log.Printf("Run K8sResourceWatcher with kind %v\n", r.kind)

	// Start our controller.
	if err := r.ctrl.Run(r.stopC); err != nil {
		return fmt.Errorf("error running controller: %s", err)
	}
	return nil
}

// Shutdown the resource watcher
func (r *K8sResourceWatcher) Shutdown() {
	// FIXME: Shutdown is not really stopping the kooper controller, althought upstream
	// the stop signal is trigerring a general shutdown followed by exit, this may change in the future!!
	log.Printf("Shutdown signal received for K8sResourceWatcher with kind %v, the controller is not being explicitly stopped, althought the whole kwatchman exits upstream so this is not a big deal\n", r.kind)
	r.stopC <- struct{}{}
}

func newK8sResourceWatcher(kind string, hand *handler.HandlerFunc, retr *retrieve.Resource) watcher.ResourceWatcher {
	// Create the controller that will refresh every 30 seconds.
	ctrl := newK8sController(kind, hand, retr)
	stopC := make(chan struct{})

	return &K8sResourceWatcher{
		kind:  kind,
		ctrl:  ctrl,
		stopC: stopC,
	}
}

func newK8sController(name string, hand *handler.HandlerFunc, retr *retrieve.Resource) controller.Controller {
	cfg := &controller.Config{
		Name:              name,
		ConcurrentWorkers: 1,
		ResyncInterval:    30 * time.Second,
	}
	return controller.New(cfg, hand, retr, nil, nil, metrics.Dummy, log.StandardLogger())
}
