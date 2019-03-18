package resources

import (
	"context"
	"testing"
)

func TestNewK8sDeploymentWatcher(t *testing.T) {
	w := NewK8sDeploymentWatcher(context.Background())
	dw := w.(*K8sResourceWatcher)

	if dw.ctx == nil {
		t.Errorf("K8sDeploymentWatcher.ctx is not set correctly %#v", dw.ctx)
	}
	if dw.informer == nil {
		t.Errorf("K8sDeploymentWatcher.informer is not set correctly %#v", dw.informer)
	}
	if dw.kind == "" {
		t.Errorf("K8sDeploymentWatcher.kind is not set correctly %#v", dw.kind)
	}
}
