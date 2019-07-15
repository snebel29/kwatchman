package registry_test

import (
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	// import of log_handler will register it once
	log_handler "github.com/snebel29/kwatchman/internal/pkg/handler/log"
	"testing"
)

func TestRegistry(t *testing.T) {
	registryName := "handler"

	handlerList, ok := registry.GetRegistry(registryName)
	if !ok {
		t.Errorf("Registry %s should have existed", registryName)
	}

	expected := 1
	if len(handlerList) != expected {
		t.Errorf("handlerList length should be %d but got %d instead", expected, len(handlerList))
	}

	// Force another register
	registry.Register(registryName, "log2", log_handler.NewLogHandler)
	expected = 2
	if len(handlerList) != expected {
		t.Errorf("handlerList length should be %d but got %d instead", expected, len(handlerList))
	}

}
