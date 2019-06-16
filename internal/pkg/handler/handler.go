package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/snebel29/kooper/operator/common"
)

type ResourcesHandlerFunc func(context.Context, Input) (Output, error)

// Input holds the input data for any handler
type Input struct {
	Evt         *common.K8sEvent
	K8sManifest []byte
	Payload     []byte //This is a free field that can hold, anything such as text, images, etc
}

// Ouput holds the output data from any handler execution
type Output struct {
	K8sManifest []byte
	Payload     []byte //This is a free field that can hold, anything such as text, images, etc
	RunNext     bool
}

type ChainOfHandlers interface {
	Run(ctx context.Context, input Input) error
}

// chainOfHandlers holds a list of ResourcesHandlerFunc that can be exexcute sequencially
type chainOfHandlers struct {
	handlers []ResourcesHandlerFunc
}

// Run will execute each handler one after the other, the handler itself is responsible to decide
// whether the next handler should be executed or not
func (c *chainOfHandlers) Run(ctx context.Context, input Input) error {
	toSend := input.K8sManifest
	payload := input.Payload

	for i, f := range c.handlers {
		output, err := f(ctx, Input{Evt: input.Evt, K8sManifest: toSend, Payload: payload})
		if err != nil {
			return errors.Wrapf(err, "The %d function failed within chainOfHandlers run()", i)
		}
		if !output.RunNext {
			break
		}
		toSend = output.K8sManifest
		payload = output.Payload
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
