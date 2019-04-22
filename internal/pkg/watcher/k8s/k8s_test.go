package k8s

import (
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"k8s.io/client-go/rest"
	"os"
	"path"
	"runtime"
	"testing"
)

var thisFilename string

func init() {
	_, t, _, _ := runtime.Caller(0)
	thisFilename = t
}

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
	kubeconfig := path.Join(path.Dir(thisFilename), "fixtures", "kubeconfig")
	_, err := getK8sClient(kubeconfig)
	if err != nil {
		t.Errorf("Failed to get client from kubeconfig %s: %v", kubeconfig, err)
	}
}

func TestNewK8sWatcher(t *testing.T) {
	kubeconfig := path.Join(path.Dir(thisFilename), "fixtures", "kubeconfig")
	c := &cli.CLIArgs{Kubeconfig: kubeconfig}
	w, err := NewK8sWatcher(c)

	if err != nil {
		t.Errorf("%s", err)
	}
	if w.opts != c {
		t.Errorf("K8sWatcher.opts is not set correctly %#v", w.opts)
	}
}

type ResourceWatcherMock struct {
	RunCalled      bool
	ShutdownCalled bool
}

func (w *ResourceWatcherMock) Run() error {
	w.RunCalled = true
	return nil
}

func (w *ResourceWatcherMock) Shutdown() {
	w.ShutdownCalled = true
}

func TestK8sWatcher(t *testing.T) {
	w := &K8sWatcher{
		opts: nil,
		k8sResources: []watcher.ResourceWatcher{
			&ResourceWatcherMock{},
			&ResourceWatcherMock{},
			&ResourceWatcherMock{},
		},
	}

	w.Run()
	for i, rwi := range w.k8sResources {
		rw := rwi.(*ResourceWatcherMock)
		if !rw.RunCalled {
			t.Errorf("ResourceWatcherMock %v Run() should have been called %#v", i, rw)
		}

	}

	w.Shutdown()
	for i, rwi := range w.k8sResources {
		rw := rwi.(*ResourceWatcherMock)
		if !rw.ShutdownCalled {
			t.Errorf("ResourceWatcherMock %v Shutdown() should have been called %#v", i, rw)
		}

	}
}
