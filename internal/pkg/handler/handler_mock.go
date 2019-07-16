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
func (h *MockHandler) Run(ctx context.Context, input Input) (Output, error) {
	h.Called = true
	h.PassedPayload = input.Payload
	h.PassedResourceKind = input.ResourceKind
	h.PassedK8sManifest = input.K8sManifest
	h.PassedEvent = input.Evt
	h.PassedContext = ctx

	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true}, nil
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
func (h *MockHandlerError) Run(ctx context.Context, input Input) (Output, error) {
	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     false}, fmt.Errorf("dummy error")
}

// NewMockHandler return a mock
func NewMockHandler() *MockHandler {
	return &MockHandler{Called: false}
}

// NewMockHandlerError return a mock that return error
func NewMockHandlerError() *MockHandlerError {
	return &MockHandlerError{Called: false}
}
