package log

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
)

type logHandler struct {
	config config.Handler
}

func init() {
	registry.Register(registry.HANDLER, "log", NewLogHandler)
}

// NewLogHandler return a log handler
func NewLogHandler(c config.Handler) handler.Handler {
	return &logHandler{
		config: c,
	}
}

// log handler logs the event, and can be used for testing, and troubleshooting
// by enriching your event logs with k8s events
func (h *logHandler) Run(ctx context.Context, evt *handler.Event) error {
	manifestToPrint := evt.K8sManifest
	_json, err := handler.PrettyPrintJSON(evt.K8sManifest)

	if err == nil {
		manifestToPrint = _json
	}

	log.Infof("%#v\n%s", evt.K8sEvt, string(evt.Payload))
	log.Debugf("%s", string(manifestToPrint))

	return nil
}
