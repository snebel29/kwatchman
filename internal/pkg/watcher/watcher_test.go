package watcher

import (
	"k8s.io/client-go/rest"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestGetK8sClient(t *testing.T) {
	// Test InClusterConfig() simulating being within k8s cluster
	// https://github.com/snebel29/kwatchman/blob/master/vendor/k8s.io/client-go/rest/config.go#L315-L345
	os.Setenv("KUBERNETES_SERVICE_HOST", "anyValue")
	os.Setenv("KUBERNETES_SERVICE_PORT", "anyValue")
	if _, err := getK8sClient(""); err == rest.ErrNotInCluster {
		t.Error("getK8sClientInCluster() should behave like living within cluster")
	}
	os.Unsetenv("KUBERNETES_SERVICE_HOST")
	os.Unsetenv("KUBERNETES_SERVICE_PORT")

	// Test getK8sClientOutCluster()
	_, thisFilename, _, _ := runtime.Caller(0)
	kubeconfig := path.Join(path.Dir(thisFilename), "fixtures", "kubeconfig")
	_, err := getK8sClient(kubeconfig)
	if err != nil {
		t.Errorf("Failed to get client from kubeconfig %s: %v", kubeconfig, err)
	}
}
