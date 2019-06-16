package handler

import (
	"context"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/handler/slack"
)

type Notifier interface {
	Notify() ResourcesHandlerFunc
}

type notifier struct {
	clusterName string
	notify      func(kind, key, clusterName, payload string) error
}

func NewSlackNotifier(clusterName string) *notifier {
	return &notifier{
		clusterName: clusterName,
		notify:      slack.MsgToSlack,
	}
}

func (n *notifier) Send(_ context.Context, evt *common.K8sEvent, payload []byte) ([]byte, bool, error) {
	err := n.notify(evt.Kind, evt.Key, n.clusterName, string(payload))
	if err != nil {
		return []byte{}, false, err
	}
	return payload, true, err
}
