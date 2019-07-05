package resources

import (
	"testing"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"github.com/snebel29/kwatchman/internal/pkg/handler/log"
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
	chainOfHandlers := handler.NewChainOfHandlers(log.NewLogHandler(config.Handler{}))
	dw := NewK8sDeploymentWatcher(ResourceWatcherArgs{
		Clientset: nil,
		Namespace: "",
		ChainOfHandlers: chainOfHandlers,
	})
	k8sIndividualResourceWatcherHelper(dw, t)
}
