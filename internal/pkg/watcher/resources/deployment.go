package resources

import (
	"context"
	apps_v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func NewK8sDeploymentWatcher(ctx context.Context) ResourceWatcher {
	// TODO: Do we need a watcher per api group to work with multiple k8s versions?
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{},
		&apps_v1.Deployment{},
		0, //Skip resync
		cache.Indexers{},
	)
	return &K8sResourceWatcher{
		ctx:      ctx,
		kind:     "Deployment",
		informer: informer,
	}
}
