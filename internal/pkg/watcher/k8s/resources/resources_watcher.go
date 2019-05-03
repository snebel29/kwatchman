package resources

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/snebel29/kooper/operator/controller"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

type K8sResourceWatcher struct {
	kind  string
	stopC chan struct{}
	ctrl  controller.Controller
}

func (r *K8sResourceWatcher) Run() error {
	log.Printf("Run K8sResourceWatcher with kind %v\n", r.kind)

	// Start our controller.
	if err := r.ctrl.Run(r.stopC); err != nil {
		return fmt.Errorf("error running controller: %s", err)
	}
	return nil
}

func (r *K8sResourceWatcher) Shutdown() {
	log.Printf("Shutdown K8sResourceWatcher with kind %v\n", r.kind)
	r.stopC <- struct{}{} //FIX: It's not really stopping the kooper controller
}

func newK8sResourceWatcher(kind string, hand *handler.HandlerFunc, retr *retrieve.Resource) watcher.ResourceWatcher {
	// Create the controller that will refresh every 30 seconds.
	ctrl := newK8sController(hand, retr)
	stopC := make(chan struct{})

	return &K8sResourceWatcher{
		kind:  kind,
		ctrl:  ctrl,
		stopC: stopC,
	}
}

func newK8sController(hand *handler.HandlerFunc, retr *retrieve.Resource) controller.Controller {
	refresh := 30 * time.Second
	logger := log.New()
	return controller.NewSequential(refresh, hand, retr, nil, logger)
}
