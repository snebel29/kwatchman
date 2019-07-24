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

// Handler interface
type Handler interface {
	Run(context.Context, *Event) error
}

// Event holds the input data for any handler
type Event struct {
	K8sEvt       *common.K8sEvent
	RunNext      bool
	ResourceKind string
	K8sManifest  []byte
	Payload      []byte //This is a free field that can hold, anything such as text, images, etc
}

// ChainOfHandlers Interface
type ChainOfHandlers interface {
	Run(context.Context, *Event) error
}

// chainOfHandlers holds a list of ResourcesHandlerFunc that can be executed sequencially
type chainOfHandlers struct {
	handlers []Handler
}

// Run will execute each handler one after the other, the handler itself is responsible to decide
// whether the next handler should be executed or not
func (c *chainOfHandlers) Run(ctx context.Context, evt *Event) error {
	for i, h := range c.handlers {
		err := h.Run(ctx, evt)
		if err != nil {
			return errors.Wrapf(err, "The %d function failed within chainOfHandlers run()", i)
		}
		if !evt.RunNext {
			break
		}
	}
	return nil
}

// NewChainOfHandlers return a ChainOfHandlers
func NewChainOfHandlers(handlers ...Handler) ChainOfHandlers {
	return &chainOfHandlers{
		handlers: handlers,
	}
}

// GetHandlerListFromConfig return list of handler objects from configuration
// their position in the list matches the defined user execution sequence
func GetHandlerListFromConfig(c *config.Config) ([]Handler, error) {
	var handlerList []Handler
	registeredHandlers, ok := registry.GetRegistry(registry.HANDLER)
	if !ok {
		return nil, errors.New("There is no handler registry available")
	}

	for _, configHandler := range c.Handlers {
		if rh, ok := registeredHandlers[configHandler.Name]; ok {
			regHandler, ok := rh.(func(config.Handler) Handler)
			if !ok {
				return nil, errors.Errorf(
					"handler %s is not of type func() Handler but %T instead", configHandler.Name, rh)
			}
			handlerList = append(handlerList, regHandler(configHandler))
		}
	}
	return handlerList, nil
}

// PrettyPrintJSON return an indented JSON
func PrettyPrintJSON(_json []byte) ([]byte, error) {
	var indented bytes.Buffer
	if err := json.Indent(&indented, _json, "", " "); err != nil {
		return nil, err
	}
	return []byte(indented.String()), nil
}
