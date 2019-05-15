package resources

import (
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"testing"
)

func k8sIndividualResourceWatcherHelper(w watcher.ResourceWatcher, t *testing.T) {
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
}

func TestK8sDeploymentWatcher(t *testing.T) {
	k8sIndividualResourceWatcherHelper(NewK8sDeploymentWatcher(nil, handler.LogHandlerFunc), t)
}
