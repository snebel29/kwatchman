package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
)

//TODO: Fix returning errors within handlers functions cause kwatch to panic!!!

// LogHandlerFunc can be used for debugging, troubleshooting and testing
func LogHandlerFunc(ctx context.Context, input Input) (Output, error) {
	_json, err := prettyPrintJSON(input.K8sManifest)
	if err == nil {
		input.K8sManifest = _json
	}

	log.Infof("%#v %s %s", input.Evt, string(input.K8sManifest), string(input.Payload))

	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true}, nil
}
