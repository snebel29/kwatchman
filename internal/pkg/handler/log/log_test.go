package log

import (
	"github.com/snebel29/kwatchman/internal/pkg/config"
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"reflect"
	"testing"
)

func TestLogHandlerFunc(t *testing.T) {

	hook := log_test.NewGlobal()

	h := NewLogHandler(config.Handler{})

	evt := &common.K8sEvent{}
	manifest := []byte("{\"a\": 1}")
	payload := []byte("payload")
	resourceKind := "Deployment"

	output, err := h.Run(
		nil,
		handler.Input{
			Evt:          evt,
			ResourceKind: resourceKind,
			K8sManifest:  manifest,
			Payload:      payload,
		})
	m := hook.Entries

	if len(m) != 1 {
		t.Errorf("There should be one entry, there is %d instead", len(m))
	}
	if err != nil {
		t.Error(err)
	}

	// TODO: Testing Ouput is common among some handler files, we could create a helper function
	if output.RunNext != true {
		t.Error("RunNext should be true")
	}
	if !reflect.DeepEqual(output.Payload, payload) {
		t.Errorf("Payload %s should match %s", string(output.Payload), string(payload))
	}
	if !reflect.DeepEqual(output.K8sManifest, manifest) {
		t.Errorf("K8sManifest %s should match %s", string(output.K8sManifest), string(manifest))
	}
}
