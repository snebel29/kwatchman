package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

func marshal(v interface{}) ([]byte, error) {
	manifest, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrapf(err, "Marshaling %T within getUnderlyingValue", v)
	}
	return manifest, nil
}

func getManifest(obj interface{}) ([]byte, error) {
	switch v := obj.(type) {
	case *appsv1.Deployment:
		return marshal(v)

	default:
		return nil, fmt.Errorf("Unknown type %T for %#v object", obj, obj)
	}
}

//This function handles Add, Update and Delete events
//the main difference between them, which is the evt.Object being <nil> and therefore causing
//the type assertion to fail, downstream Handler functions should apply different
// logic based on its evt.Kind value
func NewKooperHandlerFunction(
	chainOfHandlers handler.ChainOfHandlers,
	resourceKind string) func(context.Context, *kooper.K8sEvent) error {

	fn := func(_ context.Context, evt *kooper.K8sEvent) error {
		var err error
		var manifest []byte

		switch evt.Kind {
		case "Add", "Update":
			manifest, err = getManifest(evt.Object)
			if err != nil {
				return errors.Wrap(err, "Marshal within NewKooperHandlerFunction")
			}

		case "Delete":
			manifest = []byte{}

		default:
			return fmt.Errorf("Unknown evt.Kind %s", evt.Kind)
		}

		err = chainOfHandlers.Run(nil, handler.Input{
			Evt:          evt,
			ResourceKind: resourceKind,
			K8sManifest:  manifest,
			Payload:      []byte{},
		})

		if err != nil {
			return err
		}

		return nil
	}
	return fn
}

func NewK8sDeploymentWatcher(
	clientset kubernetes.Interface,
	chainOfHandlers handler.ChainOfHandlers) watcher.ResourceWatcher {

	resourceKind := "Deployment"
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

	fn := NewKooperHandlerFunction(chainOfHandlers, resourceKind)
	hand := &kooper_handler.HandlerFunc{
		AddFunc:    fn,
		DeleteFunc: fn,
	}
	return newK8sResourceWatcher(resourceKind, hand, retr)
}
