package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
)

//TODO: Fix returning errors within handlers functions cause kwatch to panic!!!

type logHandler struct{}

func NewLogHandler() *logHandler {
	return &logHandler{}
}

// log Run can be used for debugging, troubleshooting and testing
func (h *logHandler) Run(ctx context.Context, input Input) (Output, error) {
	manifestToPrint := input.K8sManifest
	_json, err := prettyPrintJSON(input.K8sManifest)

	if err == nil {
		manifestToPrint = _json
	}

	log.Infof("%#v\n%s", input.Evt, string(input.Payload))
	log.Debugf("%s", string(manifestToPrint))

	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true}, nil
}
