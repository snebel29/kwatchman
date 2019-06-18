package handler

import (
	"context"
	"fmt"
	"github.com/snebel29/kooper/operator/common"
)

type MockHandler struct {
	Called             bool
	PassedPayload      []byte
	PassedK8sManifest  []byte
	PassedResourceKind string
	PassedEvent        *common.K8sEvent
	PassedContext      context.Context
}

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

type MockHandlerError struct {
	Called             bool
	PassedPayload      []byte
	PassedK8sManifest  []byte
	PassedResourceKind string
	PassedEvent        *common.K8sEvent
	PassedContext      context.Context
}

func (h *MockHandlerError) Run(ctx context.Context, input Input) (Output, error) {
	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     false}, fmt.Errorf("dummy error")
}

func NewMockHandler() *MockHandler {
	return &MockHandler{Called: false}
}

func NewMockHandlerError() *MockHandlerError {
	return &MockHandlerError{Called: false}
}
