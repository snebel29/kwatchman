package watcher

import (
	"path"
	"runtime"
	"testing"
)

func TestGetK8sClient(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	kubeconfig := path.Join(path.Dir(filename), "fixtures", "kubeconfig")

	_, err := getK8sClient(kubeconfig)
	if err != nil {
		t.Errorf("Failed to get client from %s: %v", kubeconfig, err)
	}
	//TODO: Test rest.InClusterConfig() - Actually is hard to do
	//since that function uses hardcoded /var/run paths that requires
	//privileged permissions, the function can't be easily tested
}
