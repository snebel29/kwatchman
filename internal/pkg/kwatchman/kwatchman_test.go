package kwatchman

import (
	"github.com/bouk/monkey"
	"os"
	"syscall"
	"testing"
	"time"
)

type WatcherMock struct {
	RunCalled 	   bool
	ShutdownCalled bool
	ExitCalled     bool
}

func (w *WatcherMock) Run() error {
	w.RunCalled = true
	return nil
}

func (w *WatcherMock) Shutdown() error {
	w.ShutdownCalled = true
	return nil
}

func TestRun(t *testing.T) {
	watcherMock := new(WatcherMock)
	fakeExit := func(int) {
		watcherMock.ExitCalled = true
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()

	// Run the Start function
	Start(watcherMock)

	if !watcherMock.RunCalled {
		t.Error("watcher.Run() wasn't called")
	}

	shutdown <- syscall.SIGTERM
	time.Sleep(time.Millisecond * 500)

	if !watcherMock.ShutdownCalled {
		t.Error("watcher.Shutdown() wasn't called")
	}
	if !watcherMock.ExitCalled {
		t.Error("os.Exit() wasn't called during watcher.shutdown()")

	}
}