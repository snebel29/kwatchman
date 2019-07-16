package kwatchman

import (
	"github.com/bouk/monkey"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

type WatcherMock struct {
	sync.Mutex
	RunCalled      bool
	ShutdownCalled bool
	ExitCalled     bool
}

func (w *WatcherMock) SetRunCalled(flag bool) {
	w.Lock()
	defer w.Unlock()
	w.RunCalled = flag
}

func (w *WatcherMock) Run() error {
	w.SetRunCalled(true)
	return nil
}

func (w *WatcherMock) GetShutdownCalled() bool {
	w.Lock()
	defer w.Unlock()
	return w.ShutdownCalled
}

func (w *WatcherMock) SetShutdownCalled(flag bool) {
	w.Lock()
	defer w.Unlock()
	w.ShutdownCalled = flag
}

func (w *WatcherMock) Shutdown() {
	w.Lock()
	defer w.Unlock()
	w.ShutdownCalled = true
}

func (w *WatcherMock) SetExitCalled(flag bool) {
	w.Lock()
	defer w.Unlock()
	w.ExitCalled = flag
}

func TestStart(t *testing.T) {
	watcherMock := new(WatcherMock)
	fakeExit := func(int) {
		watcherMock.SetExitCalled(true)
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	// Run the Start function
	if err := Start(watcherMock); err != nil {
		t.Error(err)
	}

	if !watcherMock.RunCalled {
		t.Error("watcher.Run() wasn't called")
	}

	shutdown <- syscall.SIGTERM
	time.Sleep(time.Millisecond * 500)

	if !watcherMock.GetShutdownCalled() {
		t.Error("watcher.Shutdown() wasn't called")
	}

	if !watcherMock.ExitCalled {
		t.Error("os.Exit() wasn't called during watcher.shutdown()")
	}
}
