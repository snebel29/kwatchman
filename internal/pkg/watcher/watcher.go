package watcher

import (
	"fmt"
	"errors"
	_ "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
)

type Watcher struct {
	opts 	  *cli.CLIArgs
	clientset kubernetes.Interface
}

func NewWatcher(c *cli.CLIArgs) *Watcher {
	return &Watcher{
		opts: 	   c,
	}
}

func (w *Watcher) Run() error {
	clientset, err := w.GetK8sClient()
	if err != nil {
		return err
	}
	w.clientset = clientset
	// TODO: Run watcher loop ...
	return errors.New("Watcher.Run() should have never returned!")
}

// Returns kubernetes API clientset, depending on the context where kwatchman
// is run, InCluster vs local using Kubeconfig file for the last
func (w *Watcher) GetK8sClient() (kubernetes.Interface, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		// We may be out of k8s so try to read from kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", w.opts.Kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		errors.New(fmt.Sprintf("Can't create kubernetes client: %s", err))
	}
	return clientset, nil
}
