package handler

import (
	"context"
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

func (n *notifier) Send(ctx context.Context, input Input) (Output, error) {
	err := n.notify(input.Evt.Kind, input.Evt.Key, n.clusterName, string(input.Payload))
	if err != nil {
		return Output{
			K8sManifest: input.K8sManifest,
			Payload:     input.Payload,
			RunNext:     false}, err
	}
	return Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true}, err
}
