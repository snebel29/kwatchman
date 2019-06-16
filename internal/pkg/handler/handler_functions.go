package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kooper/operator/common"
)

//TODO: Fix returning errors within handlers functions cause kwatch to panic!!!

// LogHandlerFunc can be used for debugging, troubleshooting and testing
func LogHandlerFunc(_ context.Context, evt *common.K8sEvent, payload []byte) ([]byte, bool, error) {
	_json, err := prettyPrintJSON(payload)
	if err == nil {
		payload = _json
	}
	log.Info(string(payload))
	return nil, false, nil
}
