package resources

import (
	"context"
	log "github.com/sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

func NewK8sDeploymentWatcher(clientset kubernetes.Interface) watcher.ResourceWatcher {
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

	hand := &handler.HandlerFunc{
		AddFunc: func(_ context.Context, evt *common.K8sEvent) error {
			obj := evt.Object.(*appsv1.Deployment)
			log.Infof("%s added: %s/%s %v\n", kind, obj.Namespace, obj.Name, evt.HasSynced)
			return nil
		},
		DeleteFunc: func(_ context.Context, evt *common.K8sEvent) error {
			log.Infof("%s deleted: %s\n", kind, evt.Key)
			return nil
		},
	}
	return newK8sResourceWatcher(kind, hand, retr)
}
