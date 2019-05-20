package handler

import (
	"bytes"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kooper/operator/common"
)

type ResourcesHandlerFunc func(context.Context, *common.K8sEvent, []byte) error

func LogHandlerFunc(ctx context.Context, evt *common.K8sEvent, k8sManifest []byte) error {
	_ = ctx
	var indented bytes.Buffer
	json.Indent(&indented, k8sManifest, "", " ")
	log.Infof("Added: %v %s\n", evt.HasSynced, &indented)
	return nil
}
