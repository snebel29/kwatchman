package resources

import (
	"sync"
	"testing"
	//"sync"
	"time"
)

var runCalled bool

type KooperControllerMock struct{}

func (w *KooperControllerMock) Run(stopper <-chan struct{}) error {
	runCalled = true
	return nil
}

func TestK8sDeploymentWatcher(t *testing.T) {
	// Test that
	w := NewK8sDeploymentWatcher(nil)
	dw := w.(*K8sResourceWatcher)
	if dw.kind == "" {
		t.Errorf("kind should be != than \"\"")
	}
	if dw.ctrl == nil {
		t.Errorf("ctrl should be != nil")
	}
	if dw.stopC == nil {
		t.Errorf("stopC should be != nil")
	}

	// Test Shutdown() sends empty struct through stopC
	stopC := make(chan struct{})
	dw.stopC = stopC

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

	dw.Shutdown()
	wg.Wait()

	//Test Run() starts controller by calling ctrl.Run()
	dw.ctrl = &KooperControllerMock{}
	dw.Run()
	if runCalled == false {
		t.Error("ctrl.Run() should have been called")
	}
}
