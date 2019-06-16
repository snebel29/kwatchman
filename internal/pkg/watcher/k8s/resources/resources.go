package resources

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	kooper "github.com/snebel29/kooper/operator/common"
	kooper_handler "github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

func NewK8sDeploymentWatcher(
	clientset kubernetes.Interface,
	chainOfHandlers handler.ChainOfHandlers) watcher.ResourceWatcher {

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

	fn := func(_ context.Context, evt *kooper.K8sEvent) error {
		if obj, ok := evt.Object.(*appsv1.Deployment); ok {
			manifest, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			err = chainOfHandlers.Run(nil, handler.Input{
				Evt:         evt,
				K8sManifest: manifest,
				Payload:     []byte{},
			})
			if err != nil {
				return err
			}
		} else {
			log.Warnf("runtime.Object is not of type (*appsv1.Deployment) but %T instead", obj)
		}
		return nil
	}

	hand := &kooper_handler.HandlerFunc{
		AddFunc:    fn,
		DeleteFunc: fn,
	}
	return newK8sResourceWatcher(kind, hand, retr)
}
