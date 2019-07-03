package k8s

import (
	"fmt"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"github.com/snebel29/kwatchman/internal/pkg/watcher/k8s/resources"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

type K8sWatcher struct {
	config       *config.Config
	k8sResources []watcher.ResourceWatcher
}

func NewK8sWatcher(c *config.Config) (*K8sWatcher, error) {
	clientset, err := getK8sClient(c.CLI.Kubeconfig)
	if err != nil {
		return nil, err
	}

	// TODO: Move HandlerList and ResourcesList outside

	var handlerList []handler.Handler
	registeredHandlers, ok := registry.GetRegistry(registry.HANDLER)
	if !ok {
		return nil, fmt.Errorf("There is no handler registry available")
	}

	for _, h := range c.Handlers {
		if rh, ok := registeredHandlers[h.Name]; ok {
			regHandler, ok := rh.(handler.Handler)
			if !ok {
				return nil, fmt.Errorf(
					"handler %s is not of type handler.Handler but %T instead", h.Name, rh)
			}
			handlerList = append(handlerList, regHandler)	
		}
	}

	chainOfHandlers := handler.NewChainOfHandlers(
		handlerList...,
	)

	return &K8sWatcher{
		config: c,
		k8sResources: []watcher.ResourceWatcher{
			resources.NewK8sDeploymentWatcher(clientset, c.CLI.Namespace, chainOfHandlers),
		},
	}, nil
}

// Run start k8s controller for each k8s resource
func (w *K8sWatcher) Run() error {
	defer w.Shutdown()

	var wg sync.WaitGroup
	for _, rw := range w.k8sResources {
		wg.Add(1)
		go func(r watcher.ResourceWatcher) {
			defer wg.Done()
			// TODO: Handle errors?
			r.Run()
		}(rw)
	}

	wg.Wait()
	return nil
}

func (w *K8sWatcher) Shutdown() {
	for _, rw := range w.k8sResources {
		rw.Shutdown()
	}
}

// Returns kubernetes API clientset, depending on the context where kwatchman
// is run, InCluster vs local, kubeconfig will be used only when running out of k8s
// you can pass an empty string when running InCluster
func getK8sClient(kubeconfigFile string) (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	// If InClusterConfig() fails
	if err != nil {
		// If ErrNotInCluster then we try to get client from kubeconfig
		if err == rest.ErrNotInCluster {
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
			if err != nil {
				return nil, err
			}

			// If the error is something else we just fail
		} else {
			return nil, err
		}
	}

	// Generate new clientset from the config (either produced In or Out cluster)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Can't create kubernetes client: %s", err)
	}
	return clientset, nil
}
