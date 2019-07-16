package resources

import (
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

const (
	INGRESS = "ingress"
)

func init() {
	registry.Register(registry.RESOURCES, INGRESS, NewIngressWatcher)
}

// NewIngressWatcher return a watcher for k8s ingress
// TODO: From 1.14 extensions/v1beta1 apigroup is deprecated in favour of networking.k8s.io
// ideally we should handle this api group change transparently for users, although
// old users will eventually have to migrate over new versions of kwatchman with the new type
// and a go-client version compatible
func NewIngressWatcher(arg ResourceWatcherArgs) watcher.ResourceWatcher {

	resourceKind := INGRESS

	retr := &retrieve.Resource{
		Object: &extensions_v1beta1.Ingress{},
		ListerWatcher: &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.ExtensionsV1beta1().Ingresses(arg.Namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.LabelSelector = arg.LabelSelector
				return arg.Clientset.ExtensionsV1beta1().Ingresses(arg.Namespace).Watch(options)
			},
		},
	}

	return newK8sResourceWatcher(
		resourceKind, newResourceHandlerFunc(arg.ChainOfHandlers, resourceKind),
		retr)
}
