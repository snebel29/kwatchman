package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

const (
	DAEMONSET = "daemonset"
)

func init() {
	registry.Register(registry.RESOURCES, DAEMONSET, NewDaemonsetWatcher)
}

// NewDaemonsetWatcher return a watcher for k8s daemonset
func NewDaemonsetWatcher(arg ResourceWatcherArgs) watcher.ResourceWatcher {

	resourceKind := DAEMONSET

	retr := &retrieve.Resource{
		Object: &appsv1.DaemonSet{},
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.AppsV1().DaemonSets(arg.Namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.AppsV1().DaemonSets(arg.Namespace).Watch(options)
			},
		},
	}

	return newK8sResourceWatcher(
		resourceKind, newResourceHandlerFunc(arg.ChainOfHandlers, resourceKind),
		retr)
}
