package handler

import (
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"testing"
)

func TestPrettyPrintJSON(t *testing.T) {

	hook := log_test.NewGlobal()
	s := "{\"a\": 1}"
	manifest := []byte(s)
	LogHandlerFunc(nil, &common.K8sEvent{}, manifest)
	m := hook.LastEntry().Message

	if m != prettyPrintJSON([]byte(s)) {
		t.Errorf("%s should match %s", m, s)
	}
}
