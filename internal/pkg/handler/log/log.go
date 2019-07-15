package log

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
)

//TODO: Fix returning errors within handlers functions cause kwatch to panic!!!

type logHandler struct {
	config config.Handler
}

func init() {
	registry.Register(registry.HANDLER, "log", NewLogHandler)
}

func NewLogHandler(c config.Handler) handler.Handler {
	return &logHandler{
		config: c,
	}
}

// log Run can be used for debugging, troubleshooting and testing
func (h *logHandler) Run(ctx context.Context, input handler.Input) (handler.Output, error) {
	manifestToPrint := input.K8sManifest
	_json, err := handler.PrettyPrintJSON(input.K8sManifest)

	if err == nil {
		manifestToPrint = _json
	}

	log.Infof("%#v\n%s", input.Evt, string(input.Payload))
	log.Debugf("%s", string(manifestToPrint))

	return handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true}, nil
}
