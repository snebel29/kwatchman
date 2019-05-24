package handler

import (
	"fmt"
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

func TestDiffHandlerUseSingletonStorage(t *testing.T) {
	s1 := newStorage()
	addr1 := fmt.Sprintf("%p", s1)
	s2 := newStorage()
	addr2 := fmt.Sprintf("%p", s2)
	if addr1 != addr2 {
		t.Errorf("addr should be the same however %s != %s", addr1, addr2)
	}
}
