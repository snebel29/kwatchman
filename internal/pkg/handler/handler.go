package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
)

type Handler interface {
	Run(context.Context, Input) (Output, error)
}

// Input holds the input data for any handler
type Input struct {
	Evt          *common.K8sEvent
	ResourceKind string
	K8sManifest  []byte
	Payload      []byte //This is a free field that can hold, anything such as text, images, etc
}

// Ouput holds the output data from any handler execution
type Output struct {
	K8sManifest []byte
	Payload     []byte //This is a free field that can hold, anything such as text, images, etc
	RunNext     bool
}

type ChainOfHandlers interface {
	Run(context.Context, Input) error
}

// chainOfHandlers holds a list of ResourcesHandlerFunc that can be exexcute sequencially
type chainOfHandlers struct {
	handlers []Handler
}

// Run will execute each handler one after the other, the handler itself is responsible to decide
// whether the next handler should be executed or not
func (c *chainOfHandlers) Run(ctx context.Context, input Input) error {
	toSend := input.K8sManifest
	payload := input.Payload

	for i, h := range c.handlers {
		output, err := h.Run(ctx, Input{
			Evt:          input.Evt,
			ResourceKind: input.ResourceKind,
			K8sManifest:  toSend,
			Payload:      payload,
		})
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

func NewChainOfHandlers(handlers ...Handler) ChainOfHandlers {
	return &chainOfHandlers{
		handlers: handlers,
	}
}

// GetHandlerListFromConfig return list of handler objects from configuration
func GetHandlerListFromConfig(c *config.Config) ([]Handler, error) {
	var handlerList []Handler
	registeredHandlers, ok := registry.GetRegistry(registry.HANDLER)
	if !ok {
		return nil, errors.New("There is no handler registry available")
	}

	for _, h := range c.Handlers {
		if rh, ok := registeredHandlers[h.Name]; ok {
			regHandler, ok := rh.(Handler)
			if !ok {
				return nil, errors.Errorf(
					"handler %s is not of type handler.Handler but %T instead", h.Name, rh)
			}
			handlerList = append(handlerList, regHandler)	
		}
	}
	return handlerList, nil
}

func PrettyPrintJSON(_json []byte) ([]byte, error) {
	var indented bytes.Buffer
	if err := json.Indent(&indented, _json, "", " "); err != nil {
		return nil, err
	}
	return []byte(indented.String()), nil
}
