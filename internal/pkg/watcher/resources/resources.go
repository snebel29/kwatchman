package resources

import (
	"context"
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

func (w *K8sResourceWatcher) Run() error {
	return nil
}

func (w *K8sResourceWatcher) Shutdown() {

}

func (w *K8sResourceWatcher) HasSynced() bool {
	return true
}
