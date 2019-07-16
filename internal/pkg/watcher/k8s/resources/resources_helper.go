package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"

	kooper "github.com/snebel29/kooper/operator/common"
	kooper_handler "github.com/snebel29/kooper/operator/handler"
	"k8s.io/client-go/kubernetes"

	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

// ResourceWatcherArgs hold the arguments passed to instantiate resources watchers
type ResourceWatcherArgs struct {
	Clientset       kubernetes.Interface
	Namespace       string
	LabelSelector   string
	ChainOfHandlers handler.ChainOfHandlers
}

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

	case *appsv1.StatefulSet:
		return marshal(v)

	case *appsv1.DaemonSet:
		return marshal(v)

	case *corev1.Service:
		return marshal(v)

	case *extensions_v1beta1.Ingress:
		return marshal(v)

	default:
		return nil, fmt.Errorf("unknown type %T for %#v object", obj, obj)
	}
}

//This function handles Add, Update and Delete events
//the main difference between them, which is the evt.Object being <nil> and therefore causing
//the type assertion to fail, downstream Handler functions should apply different
// logic based on its evt.Kind value
func newKooperHandlerFunction(
	chainOfHandlers handler.ChainOfHandlers,
	resourceKind string) func(context.Context, *kooper.K8sEvent) error {

	fn := func(_ context.Context, evt *kooper.K8sEvent) error {
		var err error
		var manifest []byte

		switch evt.Kind {
		case "Add", "Update":
			manifest, err = getManifest(evt.Object)
			if err != nil {
				return errors.Wrap(err, "Marshal within newKooperHandlerFunction")
			}

		case "Delete":
			manifest = []byte{}

		default:
			return fmt.Errorf("unknown evt.Kind %s", evt.Kind)
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

func newResourceHandlerFunc(ch handler.ChainOfHandlers, resourceKind string) *kooper_handler.HandlerFunc {
	fn := newKooperHandlerFunction(ch, resourceKind)
	return &kooper_handler.HandlerFunc{
		AddFunc:    fn,
		DeleteFunc: fn,
	}
}

// GetResourcesFuncListFromConfig return list of resource objects from configuration
func GetResourcesFuncListFromConfig(c *config.Config) ([]func(ResourceWatcherArgs) watcher.ResourceWatcher, error) {
	var resourceList []func(ResourceWatcherArgs) watcher.ResourceWatcher
	registeredResources, ok := registry.GetRegistry(registry.RESOURCES)
	if !ok {
		return nil, errors.New("There is no resources registry available")
	}

	for _, configResource := range c.Resources {
		if rr, ok := registeredResources[configResource.Kind]; ok {
			regResource, ok := rr.(func(ResourceWatcherArgs) watcher.ResourceWatcher)
			if !ok {
				return nil, errors.Errorf(
					"resource %s is not of type func() watcher.ResourceWatcher but %T instead", configResource, rr)
			}
			resourceList = append(resourceList, regResource)
		}
	}
	return resourceList, nil
}

// GetResourceWatcherList return the list of configured resources
func GetResourceWatcherList(
	resourcesFuncList []func(ResourceWatcherArgs) watcher.ResourceWatcher,
	args ResourceWatcherArgs) []watcher.ResourceWatcher {

	var resourceWatcherList []watcher.ResourceWatcher
	for _, r := range resourcesFuncList {
		rr := r(args)
		resourceWatcherList = append(resourceWatcherList, rr)
	}
	return resourceWatcherList
}
