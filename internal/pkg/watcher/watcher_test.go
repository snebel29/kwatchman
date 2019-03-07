package watcher

import (
	"path"
	"runtime"
	"testing"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
)

func TestGetK8sClient(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	kubeconfig := path.Join(path.Dir(filename), "fixtures", "kubeconfig")

	c := &cli.CLIArgs{
		ClusterName: "testCluster",
		Kubeconfig:  kubeconfig,
	}

	w := NewWatcher(c)
	_, err := w.GetK8sClient()
	if err != nil {
		t.Errorf("Failed to get client from kubeconfig: %v", err)
	}
	//TODO: Test rest.InClusterConfig() - Actually is hard to do
	//since that function uses hardcoded /var/run paths that requires
	//privileged permissions, the function can't be easily tested
}
