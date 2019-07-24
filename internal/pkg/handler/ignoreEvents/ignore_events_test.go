package ignoreEvents

import (
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"reflect"
	"testing"
)

func TestIgnoreEventsHandler_Run(t *testing.T) {
	h := NewIgnoreEventsHandler(
		config.Handler{
			IgnoreEvents: []string{"Add", "Delete"},
		})

	// Non ignored event shoukd continue
	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	evt := &handler.Event{
		K8sEvt:       &common.K8sEvent{Kind: "Update"},
		RunNext:      true,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	}

	err := h.Run(nil, evt)

	if err != nil {
		t.Error(err)
	}

	if evt.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(evt.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(evt.Payload), string(payload))
	}
	if !reflect.DeepEqual(evt.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(evt.Payload), string(manifest))
	}

	// Ignored event should continue
	manifest = []byte("manifest")
	payload = []byte("payload")
	resourceKind = "Deployment"

	evt = &handler.Event{
		K8sEvt:       &common.K8sEvent{Kind: "Add"},
		RunNext:      true,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	}

	err = h.Run(nil, evt)
	if err != nil {
		t.Error(err)
	}
	if evt.RunNext == true {
		t.Error("RunNext should be false")
	}
}
