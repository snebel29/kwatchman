package handler

import (
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"testing"
)

func TestLogHandlerFunc(t *testing.T) {

	hook := log_test.NewGlobal()
	s := "{\"a\": 1}"
	manifest := []byte(s)
	LogHandlerFunc(nil, &common.K8sEvent{}, manifest)
	m := hook.LastEntry().Message

	_json, err := prettyPrintJSON([]byte(s))
	if err != nil {
		t.Error(err)
	}
	if m != string(_json) {
		t.Errorf("%s should match %s", m, s)
	}
}

func TestSendToSlackFunc(t *testing.T) {

}
