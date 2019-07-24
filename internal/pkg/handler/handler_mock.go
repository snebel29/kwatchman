package handler

import (
	"context"
	"fmt"
	"github.com/snebel29/kooper/operator/common"
)

// MockHandler call registry
type MockHandler struct {
	Called             bool
	PassedPayload      []byte
	PassedK8sManifest  []byte
	PassedResourceKind string
	PassedEvent        *common.K8sEvent
	PassedContext      context.Context
}

// Run the mock
func (h *MockHandler) Run(ctx context.Context, evt *Event) error {
	h.Called = true
	h.PassedPayload = evt.Payload
	h.PassedResourceKind = evt.ResourceKind
	h.PassedK8sManifest = evt.K8sManifest
	h.PassedEvent = evt.K8sEvt
	h.PassedContext = ctx

	evt.RunNext = true
	return nil
}

// MockHandlerError call registry
type MockHandlerError struct {
	Called             bool
	PassedPayload      []byte
	PassedK8sManifest  []byte
	PassedResourceKind string
	PassedEvent        *common.K8sEvent
	PassedContext      context.Context
}

// Run the mock
func (h *MockHandlerError) Run(ctx context.Context, evt *Event) error {
	evt.RunNext = false
	return fmt.Errorf("dummy error")
}

// NewMockHandler return a mock
func NewMockHandler() *MockHandler {
	return &MockHandler{Called: false}
}

// NewMockHandlerError return a mock that return error
func NewMockHandlerError() *MockHandlerError {
	return &MockHandlerError{Called: false}
}
