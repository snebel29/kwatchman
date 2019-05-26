package handler

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kooper/operator/common"
)

type ResourcesHandlerFunc func(context.Context, *common.K8sEvent, []byte) error

//TODO: Fix returning errors within handlers cause kwatch to panic!!!

// LogHandlerFunc can be used for debugging, troubleshooting and testing
func LogHandlerFunc(_ context.Context, evt *common.K8sEvent, k8sManifest []byte) error {
	log.Infof(string(prettyPrintJSON(k8sManifest)))
	return nil
}

func prettyPrintJSON(k8sManifest []byte) []byte {
	var indented bytes.Buffer
	json.Indent(&indented, k8sManifest, "", " ")
	return []byte(indented.String())
}
