package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/snebel29/kooper/operator/common"
)

type ResourcesHandlerFunc func(context.Context, *common.K8sEvent, []byte) (payload []byte, runNext bool, err error)

type ChainOfHandlers interface {
	Run(ctx context.Context, evt *common.K8sEvent, k8sManifest []byte) error
}

// chainOfHandlers holds a list of ResourcesHandlerFunc that can be exexcute sequencially
type chainOfHandlers struct {
	handlers []ResourcesHandlerFunc
}

// Run will run each handler one after other, the handler itself is responsible to decide
// whether the next handler should be executed or not
func (c *chainOfHandlers) Run(ctx context.Context, evt *common.K8sEvent, k8sManifest []byte) error {
	toSend := k8sManifest
	for i, f := range c.handlers {
		payload, runNext, err := f(ctx, evt, toSend)
		if err != nil {
			return errors.Wrapf(err, "The %d function failed within chainOfHandlers run()", i)
		}
		if !runNext {
			break
		}
		toSend = payload
	}
	return nil
}

func NewChainOfHandlers(handlers ...ResourcesHandlerFunc) *chainOfHandlers {
	return &chainOfHandlers{
		handlers: handlers,
	}
}

func prettyPrintJSON(_json []byte) ([]byte, error) {
	var indented bytes.Buffer
	if err := json.Indent(&indented, _json, "", " "); err != nil {
		return nil, err
	}
	return []byte(indented.String()), nil
}
