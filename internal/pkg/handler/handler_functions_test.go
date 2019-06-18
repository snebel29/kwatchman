package handler

import (
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"testing"
)

func TestLogHandlerFunc(t *testing.T) {

	hook := log_test.NewGlobal()
	manifest := []byte("{\"a\": 1}")
	LogHandlerFunc(
		nil,
		Input{
			Evt:         &common.K8sEvent{},
			K8sManifest: manifest,
			Payload:     []byte{},
		})
	m := hook.Entries

	if len(m) != 1 {
		t.Errorf("There should be one entry, there is %d instead", len(m))
	}
}
