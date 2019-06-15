package handler

import (
	"context"
	"fmt"
	"github.com/snebel29/kooper/operator/common"
	"reflect"
	"testing"
)

// This test case effectively test both LogHandlerFunc and prettyPrintJSON
func TestPrettyPrintJSON(t *testing.T) {

	arg := "{\"a\": 1}"
	expected := "{\n \"a\": 1\n}"

	returned, err := prettyPrintJSON([]byte(arg))
	if err != nil {
		t.Error(err)
	}
	if expected != string(returned) {
		t.Errorf("%s should match %s", expected, string(returned))
	}
}

type MockHandler struct {
	called        bool
	passedPayload []byte
	passedEvent   *common.K8sEvent
	passedContext context.Context
}

func (h *MockHandler) dummyHandlerFunc(ctx context.Context, evt *common.K8sEvent, payload []byte) ([]byte, bool, error) {
	h.called = true
	h.passedPayload = payload
	h.passedEvent = evt
	h.passedContext = ctx

	return payload, true, nil
}

func (h *MockHandler) dummyHandlerFuncThatReturnError(ctx context.Context, evt *common.K8sEvent, payload []byte) ([]byte, bool, error) {
	return []byte{}, false, fmt.Errorf("dummy error")
}

func NewHandler() *MockHandler {
	return &MockHandler{called: false}
}

func TestChainOfHandlers_Run(t *testing.T) {
	h1 := NewHandler()
	h2 := NewHandler()
	ch := NewChainOfHandlers(h1.dummyHandlerFunc, h2.dummyHandlerFunc, h1.dummyHandlerFuncThatReturnError)

	evt := &common.K8sEvent{}
	payload := []byte("payload")

	err := ch.Run(context.TODO(), evt, payload)
	if err == nil {
		t.Error("Last handler function should have returned an error")
	}

	if h1.called != true || h2.called != true {
		t.Errorf("handlers should have been called h1: %t h2: %t", h1.called, h2.called)
	}

	if !reflect.DeepEqual(h1.passedPayload, payload) || !reflect.DeepEqual(h2.passedPayload, payload) {
		t.Errorf("payload should have been passed h1: %s h2: %s", string(h1.passedPayload), string(h2.passedPayload))
	}

	if !reflect.DeepEqual(h1.passedEvent, evt) || !reflect.DeepEqual(h2.passedEvent, evt) {
		t.Errorf("event should have been passed h1: %#v h2: %#v", h1.passedEvent, h2.passedEvent)
	}

	if !reflect.DeepEqual(h1.passedContext, context.TODO()) || !reflect.DeepEqual(h2.passedContext, context.TODO()) {
		t.Errorf("context should have been passed h1: %#v h2: %#v", h1.passedContext, h2.passedContext)
	}
}
