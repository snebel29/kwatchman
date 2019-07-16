package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

const (
	SERVICE = "service"
)

func init() {
	registry.Register(registry.RESOURCES, SERVICE, NewServiceWatcher)
}

// NewServiceWatcher return a watcher for k8s services
func NewServiceWatcher(arg ResourceWatcherArgs) watcher.ResourceWatcher {

	resourceKind := SERVICE

	retr := &retrieve.Resource{
		Object: &corev1.Service{},
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.CoreV1().Services(arg.Namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.CoreV1().Services(arg.Namespace).Watch(options)
			},
		},
	}

	return newK8sResourceWatcher(
		resourceKind, newResourceHandlerFunc(arg.ChainOfHandlers, resourceKind),
		retr)
}
