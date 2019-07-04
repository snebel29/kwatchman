package handler_test

import (
	"fmt"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"context"
	"github.com/snebel29/kooper/operator/common"
	"reflect"
	"testing"
	"path"
	"os"
	"runtime"

	// For the handlers to be registered
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/diff"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/slack"
	_ "github.com/snebel29/kwatchman/internal/pkg/handler/log"
)

var thisFilename string

func init() {
	_, t, _, _ := runtime.Caller(0)
	thisFilename = t
}

// This test case effectively test both LogHandlerFunc and prettyPrintJSON
func TestPrettyPrintJSON(t *testing.T) {

	arg := "{\"a\": 1}"
	expected := "{\n \"a\": 1\n}"

	returned, err := handler.PrettyPrintJSON([]byte(arg))
	if err != nil {
		t.Error(err)
	}
	if expected != string(returned) {
		t.Errorf("%s should match %s", expected, string(returned))
	}
}

func TestChainOfHandlers_Run(t *testing.T) {
	h1 := handler.NewMockHandler()
	h2 := handler.NewMockHandler()
	h3 := handler.NewMockHandlerError()

	ch := handler.NewChainOfHandlers(h1, h2, h3)

	evt := &common.K8sEvent{}
	manifest := []byte("manifest")
	payload := []byte("payload")
	resourceKind := "Deployment"

	err := ch.Run(context.TODO(), handler.Input{
		Evt:          evt,
		ResourceKind: resourceKind,
		K8sManifest:  manifest,
		Payload:      payload,
	})

	if err == nil {
		t.Error("Last handler function should have returned an error")
	}

	if h1.Called != true || h2.Called != true {
		t.Errorf("handlers should have been called h1: %t h2: %t", h1.Called, h2.Called)
	}

	if h1.PassedResourceKind != resourceKind ||
		h2.PassedResourceKind != resourceKind {
		t.Errorf("%s should match %s", h1.PassedResourceKind, resourceKind)
	}

	if !reflect.DeepEqual(h1.PassedPayload, payload) || !reflect.DeepEqual(h2.PassedPayload, payload) {
		t.Errorf("payload should have been passed h1: %s h2: %s", string(h1.PassedPayload), string(h2.PassedPayload))
	}

	if !reflect.DeepEqual(h1.PassedEvent, evt) || !reflect.DeepEqual(h2.PassedEvent, evt) {
		t.Errorf("event should have been passed h1: %#v h2: %#v", h1.PassedEvent, h2.PassedEvent)
	}

	if !reflect.DeepEqual(h1.PassedContext, context.TODO()) || !reflect.DeepEqual(h2.PassedContext, context.TODO()) {
		t.Errorf("context should have been passed h1: %#v h2: %#v", h1.PassedContext, h2.PassedContext)
	}
}

func TestGetHandlerListFromConfig(t *testing.T) {
	configFile := path.Join(path.Dir(thisFilename), "fixtures", "config.toml")
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--config=%s", configFile),
	}
	

	conf, err := config.NewConfig()
	if err != nil {
		t.Error("The config should have been parsed without errors")
	}
	handlerList, err := handler.GetHandlerListFromConfig(conf)
	if err != nil {
		t.Error("The handler list should have been returned without errors")
	}
	if len(handlerList) != 3 {
		t.Errorf("handlerList should have 3 handlers, have %d instead", len(handlerList))
	}
}
