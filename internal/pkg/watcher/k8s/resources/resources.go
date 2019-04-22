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
	hand  handler.Handler
	retr  retrieve.Retriever
	stopC chan struct{}
}

func (r *K8sResourceWatcher) Run() error {
	log.Printf("Run K8sResourceWatcher with kind %v\n", r.kind)

	// Create the controller that will refresh every 30 seconds.
	ctrl := controller.NewSequential(30*time.Second, r.hand, r.retr, nil, log.New())

	// Start our controller.
	if err := ctrl.Run(r.stopC); err != nil {
		return fmt.Errorf("error running controller: %s", err)
	}
	return nil
}

func (r *K8sResourceWatcher) Shutdown() {
	log.Printf("Shutdown K8sResourceWatcher with kind %v\n", r.kind)
	r.stopC <- struct{}{} //FIX: It's not really stopping the kooper controller
}

func NewK8sDeploymentWatcher(clientset kubernetes.Interface) watcher.ResourceWatcher {

	kind := "Deployment"

	// Our domain logic that will print every add/sync/update and delete event we .
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

	// Create our retriever so the controller knows how to get/listen for pod events.
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

	stopC := make(chan struct{})

	return &K8sResourceWatcher{
		kind:  kind,
		hand:  hand,
		retr:  retr,
		stopC: stopC,
	}
}
