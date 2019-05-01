package resources

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kooper/operator/controller"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

type K8sResourceWatcher struct {
	kind  string
	stopC chan struct{}
	ctrl  controller.Controller
}

func (r *K8sResourceWatcher) Run() error {
	log.Printf("Run K8sResourceWatcher with kind %v\n", r.kind)

	// Start our controller.
	if err := r.ctrl.Run(r.stopC); err != nil {
		return fmt.Errorf("error running controller: %s", err)
	}
	return nil
}

func (r *K8sResourceWatcher) Shutdown() {
	log.Printf("Shutdown K8sResourceWatcher with kind %v\n", r.kind)
	r.stopC <- struct{}{} //FIX: It's not really stopping the kooper controller
}

func newK8sResourceWatcher(kind string, hand *handler.HandlerFunc, retr *retrieve.Resource) watcher.ResourceWatcher {
	// Create the controller that will refresh every 30 seconds.
	ctrl := newK8sController(hand, retr)
	stopC := make(chan struct{})

	return &K8sResourceWatcher{
		kind:  kind,
		ctrl:  ctrl,
		stopC: stopC,
	}
}

func newK8sController(hand *handler.HandlerFunc, retr *retrieve.Resource) controller.Controller {
	refresh := 30 * time.Second
	logger := log.New()
	return controller.NewSequential(refresh, hand, retr, nil, logger)
}

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
