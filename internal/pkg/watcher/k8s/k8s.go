package k8s

import (
	"fmt"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"github.com/snebel29/kwatchman/internal/pkg/watcher/k8s/resources"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

type K8sWatcher struct {
	opts         *cli.CLIArgs
	k8sResources []watcher.ResourceWatcher
}

func NewK8sWatcher(c *cli.CLIArgs) (*K8sWatcher, error) {
	clientset, err := getK8sClient(c.Kubeconfig)
	if err != nil {
		return nil, err
	}

	chainOfHandlers := handler.NewChainOfHandlers(
		handler.DiffFunc,
		handler.NewSlackNotifier(c.ClusterName).Send,
	)

	return &K8sWatcher{
		opts: c,
		// TODO: Make resources configurable by user
		k8sResources: []watcher.ResourceWatcher{
			resources.NewK8sDeploymentWatcher(clientset, chainOfHandlers),
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
