package watcher

import (
	"errors"
	"fmt"
	_ "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Watcher interface {
	Run() error
	Shutdown()
}

type K8sWatcher struct {
	opts      *cli.CLIArgs
	clientset kubernetes.Interface
}

func NewK8sWatcher(c *cli.CLIArgs) (*K8sWatcher, error) {
	clientset, err := getK8sClient(c.Kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset = clientset
	return &K8sWatcher{
		opts:      c,
		clientset: clientset,
	}, nil
}

func (w *K8sWatcher) Run() error {
	// TODO: Run watcher loop ...
	return errors.New("Watcher.Run() should have never returned!")
}

func (w *K8sWatcher) Shutdown() {
	// TODO: Handle shutdown gracefully
}

// Returns kubernetes API clientset, depending on the context where kwatchman
// is run, InCluster vs local using Kubeconfig file for the last
func getK8sClient(kubeconfigFile string) (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		// We may be out of k8s so try to read from kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Can't create kubernetes client: %s", err))
	}
	return clientset, nil
}
