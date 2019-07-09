package k8s

import (
	"fmt"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"github.com/snebel29/kwatchman/internal/pkg/watcher/k8s/resources"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"

	//Register handlers to be available for configuration
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/log"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/diff"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/slack"
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

	handlerList, err := handler.GetHandlerListFromConfig(c)
	if err != nil {
		return nil, err
	}

	chainOfHandlers := handler.NewChainOfHandlers(handlerList...)

	resourcesFuncList, err := resources.GetResourcesFuncListFromConfig(c)
	if err != nil {
		return nil, err
	}

	return &K8sWatcher{
		config: c,
		k8sResources: resources.GetResourceWatcherList(
			resourcesFuncList,
			resources.ResourceWatcherArgs{
				Clientset: clientset,
				Namespace: c.CLI.Namespace,
				ChainOfHandlers: chainOfHandlers,
		}),
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
