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
	config config.Handler
}

// NewSlackHandler return the slack handler
func NewIgnoreEventsHandler(c config.Handler) handler.Handler {
	return &IgnoreEventsHandler{config: c}
}

// IgnoreEvents handler stop chain execution if the event kind is on the configured list
func (h *IgnoreEventsHandler) Run(ctx context.Context, evt *handler.Event) error {
	for _, event := range h.config.IgnoreEvents {
		if evt.K8sEvt.Kind == event {
			evt.RunNext = false
			return nil
		}
	}
	evt.RunNext = true
	return nil
}
