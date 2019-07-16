package k8s

import (
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	// We need handler/log init() registeting the handler for testing
	"errors"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/log"
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
	if err := os.Setenv("KUBERNETES_SERVICE_HOST", "anyValue"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("KUBERNETES_SERVICE_PORT", "anyValue"); err != nil {
		t.Error(err)
	}
	if _, err := getK8sClient(""); err == rest.ErrNotInCluster {
		t.Error("getK8sClientInCluster() should behave like living within cluster")
	}
	if err := os.Unsetenv("KUBERNETES_SERVICE_HOST"); err != nil {
		t.Error(err)
	}
	if err := os.Unsetenv("KUBERNETES_SERVICE_PORT"); err != nil {
		t.Error(err)
	}

	// Test getK8sClientOutCluster()
	kubeconfig := path.Join(path.Dir(thisFilename), "fixtures", "kubeconfig")
	_, err := getK8sClient(kubeconfig)
	if err != nil {
		t.Errorf("Failed to get client from kubeconfig %s: %v", kubeconfig, err)
	}
}

func TestNewK8sWatcher(t *testing.T) {
	kubeconfig := path.Join(path.Dir(thisFilename), "fixtures", "kubeconfig")
	handlerToRegister := "log"
	h := config.Handlers{{Name: handlerToRegister}}
	r := config.Resources{
		{Kind: "deployment", Policies: []string{}},
	}
	c := &cli.Args{
		Namespace:  "namespace",
		Kubeconfig: kubeconfig,
		ConfigFile: "",
	}
	conf := &config.Config{
		Handlers:  h,
		Resources: r,
		CLI:       c,
	}

	w, err := NewK8sWatcher(conf)

	if err != nil {
		t.Errorf("%s getting NewK8sWatcher", err)
	}
	if w.config != conf {
		t.Errorf("K8sWatcher.config is not set correctly %#v", w.config)
	}
	if len(w.k8sResources) != 1 {
		t.Errorf("There should be 1 resource, but there is %d instead", len(w.k8sResources))
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

type ResourceWatcherWithErrorMock struct {
	RunCalled      bool
	ShutdownCalled bool
}

func (w *ResourceWatcherWithErrorMock) Run() error {
	w.RunCalled = true
	return errors.New("simulated error")
}

func (w *ResourceWatcherWithErrorMock) Shutdown() {
	w.ShutdownCalled = true
}

func TestK8sWatcherRunAndShutdownNormally(t *testing.T) {
	w := &K8sWatcher{
		config: nil,
		k8sResources: []watcher.ResourceWatcher{
			&ResourceWatcherMock{},
			&ResourceWatcherMock{},
			&ResourceWatcherMock{},
		},
	}

	err := w.Run()
	if err != nil {
		t.Error(err)
	}

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

func TestK8sWatcherRunAndFailWithErrors(t *testing.T) {
	w := &K8sWatcher{
		config: nil,
		k8sResources: []watcher.ResourceWatcher{
			&ResourceWatcherMock{},
			&ResourceWatcherWithErrorMock{},
			&ResourceWatcherMock{},
		},
	}

	err := w.Run()
	if err == nil {
		t.Errorf("An error should have being returned %s", err)
	}
}
