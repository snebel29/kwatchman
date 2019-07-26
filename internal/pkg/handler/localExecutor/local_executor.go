package localExecutor

import (
	"context"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
)

func init() {
	registry.Register(registry.HANDLER, "localExecutor", NewLocalExecutorHandler)
}

type localExecutorHandler struct {
	config config.Handler
}

// NewSlackHandler return the slack handler
func NewLocalExecutorHandler(c config.Handler) handler.Handler {
	return &localExecutorHandler{config: c}
}

// IgnoreEvents handler stop chain execution if the event kind is on the configured list
func (h *localExecutorHandler) Run(ctx context.Context, evt *handler.Event) error {
	if h.config.TimeoutSeconds == 0 {
		h.config.TimeoutSeconds = 1
	}
	// stderr and exit code for error
	// Send the marshaled event struct
	// reads the event struct back and update the current event
	// times out, by default after 1 second
	return nil
}
