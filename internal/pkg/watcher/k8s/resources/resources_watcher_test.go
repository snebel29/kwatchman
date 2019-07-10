package resources

import (
	"github.com/snebel29/kooper/operator/handler"
	"github.com/snebel29/kooper/operator/retrieve"
	"sync"
	"testing"
	"time"
)

var runCalled bool

type KooperControllerMock struct{}

func (w *KooperControllerMock) Run(stopper <-chan struct{}) error {
	runCalled = true
	return nil
}

func TestK8sResourceWatcher(t *testing.T) {
	kind := "foo"
	w := newK8sResourceWatcher(kind, &handler.HandlerFunc{}, &retrieve.Resource{})
	rw := w.(*K8sResourceWatcher)

	if rw.kind == "" {
		t.Errorf("kind should be != than \"\"")
	}
	if rw.ctrl == nil {
		t.Errorf("ctrl should be != nil")
	}
	if rw.stopC == nil {
		t.Errorf("stopC should be != nil")
	}

	// Test Shutdown() sends empty struct through stopC
	stopC := make(chan struct{})
	rw.stopC = stopC

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		select {
		case <-stopC:
		case <-time.After(1 * time.Second):
			t.Errorf("Timed out while waiting for stopC signal")
		}
	}()

	rw.Shutdown()
	wg.Wait()

	//Test Run() starts controller by calling ctrl.Run()
	rw.ctrl = &KooperControllerMock{}
	if err := rw.Run(); err != nil {
		t.Error(err)
	}
	if runCalled == false {
		t.Error("ctrl.Run() should have been called")
	}
}
