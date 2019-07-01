package handler_test

import (
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	log_handler "github.com/snebel29/kwatchman/internal/pkg/handler/log"
	"testing"
)

func TestRegistry(t *testing.T) {
	handler.Register("log", log_handler.NewLogHandler())
	handlerlist := handler.GetRegistry()
	expected := 1
	if len(handlerlist) != expected {
		t.Errorf("handlerList length should be %d but got %d instead", expected, len(handlerlist))
	}
}
