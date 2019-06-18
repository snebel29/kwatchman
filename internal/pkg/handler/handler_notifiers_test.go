package handler

import (
	"errors"
	"github.com/snebel29/kooper/operator/common"
	"reflect"
	"testing"
)

type MockNotifier struct {
	called      bool
	returnError bool
}

func (m *MockNotifier) notify(kind, key, clusterName, resourceKind, payload string) error {
	m.called = true
	if m.returnError {
		return errors.New("fake error")
	}
	return nil
}

func TestSlackNotifier(t *testing.T) {
	s := NewSlackNotifier("clusterName")
	m := &MockNotifier{}
	s.notify = m.notify
	p := []byte{}
	output, err := s.Send(nil, Input{
		Evt:         &common.K8sEvent{},
		K8sManifest: []byte{},
		Payload:     p,
	})
	if err != nil {
		t.Error("No error should have been returned")
	}
	if !reflect.DeepEqual(output.Payload, p) {
		t.Error("Returned payload should match with sent one")
	}
	if output.RunNext != true {
		t.Error("Successfull execution should runNext")
	}
	if m.called == false {
		t.Error("Notify function should have been called")
	}
}
