package resources

import (
	"context"
	apps_v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewK8sDeploymentWatcher(ctx context.Context) ResourceWatcher {
	// TODO: Do we need a watcher per api group to work with multiple k8s versions?

	var kind string = "Deployment"
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{},
		&apps_v1.Deployment{},
		0, //Skip resync
		cache.Indexers{},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				event := &Event{
					key:       key,
					action:    "CREATE",
					namespace: obj.(apps_v1.Deployment).ObjectMeta.Namespace,
					kind:      kind,
				}
				queue.Add(event)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(oldObj)
			if err == nil {
				event := &Event{
					key:       key,
					action:    "UPDATE",
					namespace: newObj.(apps_v1.Deployment).ObjectMeta.Namespace,
					kind:      kind,
				}
				queue.Add(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				event := &Event{
					key:       key,
					action:    "DELETE",
					namespace: obj.(apps_v1.Deployment).ObjectMeta.Namespace,
					kind:      kind,
				}
				queue.Add(event)
			}
		},
	})
	return &K8sResourceWatcher{
		ctx:      ctx,
		kind:     kind,
		informer: informer,
		queue:    queue,
	}
}
