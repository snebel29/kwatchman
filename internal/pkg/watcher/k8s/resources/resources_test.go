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


func TestNewResourceWatcher(t *testing.T) {

	resourcesFactoryToTest := []func(ResourceWatcherArgs) watcher.ResourceWatcher{
		NewDeploymentWatcher,
		NewStatefulsetWatcher,
		NewDaemonsetWatcher,
		NewServiceWatcher,
		NewIngressWatcher,
	}

	chainOfHandlers := handler.NewChainOfHandlers(log.NewLogHandler(config.Handler{}))
	rwa := ResourceWatcherArgs{
		Clientset: nil,
		Namespace: "",
		ChainOfHandlers: chainOfHandlers,
	}

	for _, r := range resourcesFactoryToTest {
		k8sIndividualResourceWatcherHelper(r(rwa), t)
	}
}
