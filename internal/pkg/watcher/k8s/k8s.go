package k8s

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"github.com/snebel29/kwatchman/internal/pkg/watcher/k8s/resources"
	"sync"

	// We need to register cloud auth providers
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	//Register the following handlers to be available for configuration
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/diff"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/log"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/slack"
)

type Watcher struct {
	config       *config.Config
	k8sResources []watcher.ResourceWatcher
}

// NewK8sWatcher parses the config and maps handlers and
// resources from configuration, then return the k8sWatcher
func NewK8sWatcher(c *config.Config) (*Watcher, error) {
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

	return &Watcher{
		config: c,
		k8sResources: resources.GetResourceWatcherList(
			resourcesFuncList,
			resources.ResourceWatcherArgs{
				Clientset:       clientset,
				Namespace:       c.CLI.Namespace,
				LabelSelector:   c.CLI.LabelSelector,
				ChainOfHandlers: chainOfHandlers,
			}),
	}, nil
}

// Run start k8s controller for each k8s resource
func (w *Watcher) Run() error {
	// Mo matter what, either an error or a legitime shutdown returning nil,
	// at the end we shutdown the watcher with all its ResourceWatchers
	defer w.Shutdown()

	var wg sync.WaitGroup

	// errC will block until either all controllers finish or any of them return an error
	errC := make(chan error, 1)

	// Run each controller on its own goroutine and block until finish or failed
	for _, rw := range w.k8sResources {
		wg.Add(1)
		// For extra safety we pass resource watcher and channel as parameters
		// although the wait group can be safely "closurized"
		go func(r watcher.ResourceWatcher, errC chan<- error) {
			defer wg.Done()
			if err := r.Run(); err != nil {
				errC <- errors.Wrap(err, "K8sWatcher Run()")
			}
		}(rw, errC)
	}

	// Effectively blocks until all controllers finish its execution, this is a controlled shutdown
	go func() {
		wg.Wait()
		errC <- nil
	}()

	// Return either an error or nil
	return <-errC
}

func (w *Watcher) Shutdown() {
	for _, rw := range w.k8sResources {
		rw.Shutdown()
	}
}

// Returns kubernetes API clientset, depending on the context where kwatchman
// is run, InCluster vs local, kubeconfig will be used only when running out of k8s
// you can pass an empty string when running InCluster
func getK8sClient(kubeconfigFile string) (kubernetes.Interface, error) {
	var conf *rest.Config
	conf, err := rest.InClusterConfig()

	// If InClusterConfig() fails
	if err != nil {
		// If ErrNotInCluster then we try to get client from kubeconfig
		if err == rest.ErrNotInCluster {
			conf, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
			if err != nil {
				return nil, err
			}

			// If the error is something else we just fail
		} else {
			return nil, err
		}
	}

	// Generate new clientset from the config (either produced In or Out cluster)
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("can't create kubernetes client: %s", err)
	}
	return clientset, nil
}
