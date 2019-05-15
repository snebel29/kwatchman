package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kooper/operator/common"
)

type ResourcesHandlerFunc func(context.Context, *common.K8sEvent, string) error

func LogHandlerFunc(ctx context.Context, evt *common.K8sEvent, k8sManifest string) error {
	_ = ctx
	log.Infof("Added: %v %s\n", evt.HasSynced, k8sManifest)
	return nil
}
