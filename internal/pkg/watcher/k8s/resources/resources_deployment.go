package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	kooper_handler "github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

const (
	DEPLOYMENT = "deployment"
)

func init() {
	registry.Register(registry.RESOURCES, DEPLOYMENT, NewDeploymentWatcher)
}

// NewDeploymentWatcher return a watcher for k8s deployments
func NewDeploymentWatcher(arg ResourceWatcherArgs) watcher.ResourceWatcher {

	retr := &retrieve.Resource{
		Object: &appsv1.Deployment{},
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return arg.Clientset.AppsV1().Deployments(arg.Namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return arg.Clientset.AppsV1().Deployments(arg.Namespace).Watch(options)
			},
		},
	}

	fn := newKooperHandlerFunction(arg.ChainOfHandlers, DEPLOYMENT)
	hand := &kooper_handler.HandlerFunc{
		AddFunc:    fn,
		DeleteFunc: fn,
	}
	return newK8sResourceWatcher(DEPLOYMENT, hand, retr)
}
