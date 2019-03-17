package watcher

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// TODO: Allow to configure at runtime
	DeploymentEnabled bool = true
)

type Watcher interface {
	Run() error
	Shutdown()
}

type K8sWatcher struct {
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *cli.CLIArgs
	clientset kubernetes.Interface
}

func NewK8sWatcher(c *cli.CLIArgs) (*K8sWatcher, error) {
	clientset, err := getK8sClient(c.Kubeconfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &K8sWatcher{
		ctx:       ctx,
		cancel:    cancel,
		opts:      c,
		clientset: clientset,
	}, nil
}

func (w *K8sWatcher) Run() error {
	defer w.Shutdown()

	if DeploymentEnabled {
		// TODO: go func() watchDeployment
	}
	return nil
}

func (w *K8sWatcher) Shutdown() {
	w.cancel()
}

// Returns kubernetes API clientset, depending on the context where kwatchman
// is run, InCluster vs local, kubeconfig will be used only when running out of k8s
// you can pass an empty string when running InCluster
func getK8sClient(kubeconfigFile string) (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		// If ErrNotInCluster then we try to get client from kubeconfig
		if err == rest.ErrNotInCluster {
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't create kubernetes client: %s", err))
	}
	return clientset, nil
}
