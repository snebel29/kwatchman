package log

import (
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"reflect"
	"testing"
)

func TestLogHandlerFunc(t *testing.T) {

	hook := log_test.NewGlobal()
	h := NewLogHandler(config.Handler{})

	manifest := []byte("{\"a\": 1}")
	payload := []byte("payload")
	resourceKind := "Deployment"

	evt := &handler.Event{
		K8sEvt:       &common.K8sEvent{},
		RunNext:      true,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	}

	err := h.Run(nil, evt)
	m := hook.Entries

	if len(m) != 1 {
		t.Errorf("There should be one entry, there is %d instead", len(m))
	}
	if err != nil {
		t.Error(err)
	}

	// TODO: Testing Ouput is common among some handler files, we could create a helper function
	if evt.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(evt.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(evt.Payload), string(payload))
	}
	if !reflect.DeepEqual(evt.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(evt.K8sManifest), string(manifest))
	}
}
