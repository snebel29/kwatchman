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
	log.Infof(prettyPrintJSON(k8sManifest))
	return nil
}

func prettyPrintJSON(k8sManifest []byte) string {
	var indented bytes.Buffer
	json.Indent(&indented, k8sManifest, "", " ")
	return indented.String()
}
