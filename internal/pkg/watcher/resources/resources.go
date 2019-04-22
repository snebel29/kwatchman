package resources

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"

	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kooper/operator/controller"
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
)

type ResourceWatcher interface {
	Run() error
	Shutdown()
}

type K8sResourceWatcher struct {
	kind string
	hand handler.Handler
	retr retrieve.Retriever
}

func (r *K8sResourceWatcher) Run() error {
	fmt.Printf("Run K8sResourceWatcher with kind %v\n", r.kind)

	// Create the controller that will refresh every 30 seconds.
	ctrl := controller.NewSequential(30*time.Second, r.hand, r.retr, nil, log.New())

	// Start our controller.
	stopC := make(chan struct{})
	if err := ctrl.Run(stopC); err != nil {
		return fmt.Errorf("error running controller: %s", err)
	}
	return nil
}

func (r *K8sResourceWatcher) Shutdown() {
	fmt.Printf("Shutdown K8sResourceWatcher with kind %v\n", r.kind)
}

func NewK8sDeploymentWatcher(clientset kubernetes.Interface) ResourceWatcher {

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
			pod := evt.Object.(*appsv1.Deployment)
			log.Infof("Deployment added: %s/%s %v\n", pod.Namespace, pod.Name, evt.HasSynced)
			return nil
		},
		DeleteFunc: func(_ context.Context, evt *common.K8sEvent) error {
			fmt.Printf("Dep;loyment deleted: %s\n", evt.Key)
			return nil
		},
	}

	return &K8sResourceWatcher{
		kind: "Deployment",
		hand: hand,
		retr: retr,
	}
}
