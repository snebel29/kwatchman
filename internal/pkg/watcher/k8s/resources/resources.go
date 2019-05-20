package resources

import (
	"context"
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/common"
	kooper_handler "github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

func NewK8sDeploymentWatcher(
	clientset kubernetes.Interface,
	f handler.ResourcesHandlerFunc) watcher.ResourceWatcher {

	kind := "Deployment"
	retr := &retrieve.Resource{
		Object: &appsv1.Deployment{},
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientset.AppsV1().Deployments("").List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return clientset.AppsV1().Deployments("").Watch(options)
			},
		},
	}

	fn := func(_ context.Context, evt *common.K8sEvent) error {
		obj := evt.Object.(*appsv1.Deployment)
		manifest, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		if err := f(nil, evt, manifest); err != nil {
			return err
		}
		return nil
	}

	hand := &kooper_handler.HandlerFunc{
		AddFunc:    fn,
		DeleteFunc: fn,
	}
	return newK8sResourceWatcher(kind, hand, retr)
}
