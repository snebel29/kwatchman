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
	notify func(kind, key, payload string) error
}

func NewSlackNotifier() *notifier {
	return &notifier{notify: slack.MsgToSlack}
}

func (n *notifier) Send(_ context.Context, evt *common.K8sEvent, payload []byte) ([]byte, bool, error) {
	err := n.notify(evt.Kind, evt.Key, string(payload))
	if err != nil {
		return []byte{}, false, err
	}
	return payload, true, err
}
