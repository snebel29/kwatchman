package ignoreEvents

import (
	"context"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
)

func init() {
	registry.Register(registry.HANDLER, "ignoreEvents", NewIgnoreEventsHandler)
}

type IgnoreEventsHandler struct {
	config      config.Handler
}

// NewSlackHandler return the slack handler
func NewIgnoreEventsHandler(c config.Handler) handler.Handler {
	return &IgnoreEventsHandler{config: c}
}

func (h *IgnoreEventsHandler) noErrorNoRunNext() (handler.Output, error) {
	return handler.Output{RunNext: false}, nil
}

// IgnoreEvents handler stop chain execution if the event kind is on the configured list
func (h *IgnoreEventsHandler) Run(ctx context.Context, input handler.Input) (handler.Output, error) {
	for _, event := range h.config.IgnoreEvents {
		if input.Evt.Kind == event {
			return h.noErrorNoRunNext()
		}
	}
	return 	handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true,
	}, nil
}
