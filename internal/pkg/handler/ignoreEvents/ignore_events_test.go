package ignoreEvents

import (
	"reflect"
	"testing"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
)

func TestIgnoreEventsHandler_Run(t *testing.T) {
	h := NewIgnoreEventsHandler(
		config.Handler{
			IgnoreEvents: []string{"Add", "Delete"},
		})

	// Non ignored event shoukd continue
	evt          := &common.K8sEvent{Kind: "Update"}
	manifest     := []byte("manifest")
	payload      := []byte("payload")
	resourceKind := "Deployment"

	output, err := h.Run(nil, handler.Input{
		Evt:          evt,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	})

	if err != nil {
		t.Error(err)
	}

	if output.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(output.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(output.Payload), string(payload))
	}
	if !reflect.DeepEqual(output.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(output.Payload), string(manifest))
	}

	// Ignored event should continue
	evt          = &common.K8sEvent{Kind: "Add"}
	manifest     = []byte("manifest")
	payload      = []byte("payload")
	resourceKind = "Deployment"

	output, err = h.Run(nil, handler.Input{
		Evt:          evt,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	})
	if err != nil {
		t.Error(err)
	}
	if output.RunNext == true {
		t.Error("RunNext should be false")
	}
}

