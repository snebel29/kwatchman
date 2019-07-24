package diff

import (
	"bytes"
	"context"
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"io/ioutil"
	"testing"
)

func TestDiffHandler(t *testing.T) {
	hook := log_test.NewGlobal()

	// Fake JSON struct must have some common fields with k8sObject struct
	// In order to unmarshal the differences

	h1 := NewDiffHandler(config.Handler{})
	key := "key1"

	// Add event with HasSynced == false should stop and return no error
	evt := &handler.Event{
		K8sEvt: &common.K8sEvent{
			Key:       "anyKey",
			HasSynced: false,
			Kind:      "Add",
			Object:    nil,
		},
		RunNext:     true,
		K8sManifest: []byte("{\"kind\": \"fakeKind\"}\n"),
		Payload:     []byte{},
	}

	err := h1.Run(context.TODO(), evt)
	if err != nil {
		t.Errorf("NO error should have been raised, got %s instead", err)
	}
	if evt.RunNext == true {
		t.Errorf("Add event with HasSynced == false should stop evt: %#v", evt)
	}

	// A new key is added, no difference, no error, and nextRun == true should be expected

	evt = &handler.Event{
		K8sEvt: &common.K8sEvent{
			Key:       key,
			HasSynced: true,
			Kind:      "Add",
			Object:    nil,
		},
		RunNext:     true,
		K8sManifest: []byte("{\"kind\": \"fakeKind\"}\n"),
		Payload:     []byte{},
	}

	err = h1.Run(context.TODO(), evt)

	diff := evt.Payload
	if len(diff) != 0 {
		t.Error("There should be no diff because is a new event")
	}

	if evt.RunNext != true {
		t.Error("nextRun should be true because this is an Add event")
	}

	if err != nil {
		t.Error("No error should have ocurred")
	}

	if hook.LastEntry() != nil {
		t.Errorf("Logging lastEntry should be nil")
	}

	// The same event (we use the previous one) with the same kind no difference should be returned
	// and the handler should stop
	evt.K8sEvt.Kind = "Update"
	err = h1.Run(context.TODO(), evt)

	diff = evt.Payload
	if len(diff) > 0 {
		t.Error("there should be NO difference")
	}

	if evt.RunNext != false {
		t.Errorf("RunNext should be false because NO difference was found evt: %#v", evt)
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	// The same key with diffrent kind is Updated on the same handler
	// now a difference should be returned
	evt = &handler.Event{
		K8sEvt: &common.K8sEvent{
			Key:       key,
			HasSynced: true,
			Kind:      "Update",
			Object:    nil,
		},
		RunNext:     true,
		K8sManifest: []byte("{\"kind\": \"fakeKindDifferentThanPrevious\"}\n"),
		Payload:     []byte{},
	}
	err = h1.Run(context.TODO(), evt)

	diff = evt.Payload
	if len(diff) < 1 {
		t.Error("there should be some difference")
	}

	if evt.RunNext != true {
		t.Error("nextRun should be true because difference was found")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	// Manifest should be deleted from storage and next handlers not to be trigger
	evt = &handler.Event{
		K8sEvt: &common.K8sEvent{
			Key:       key,
			HasSynced: true,
			Kind:      "Delete",
			Object:    nil,
		},
		RunNext:     true,
		K8sManifest: []byte{},
		Payload:     []byte{},
	}

	err = h1.Run(context.TODO(), evt)

	diff = evt.Payload
	if len(diff) != 0 {
		t.Error("There should be no difference")
	}

	if evt.RunNext != true {
		t.Error("nextRun should be true because deletes are notified")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}
}

func TestUpdateWithWrongHandlerCommandShouldReturnError(t *testing.T) {
	h := &diffHandler{
		config:      config.Handler{},
		storage:     newStorage(),
		diffCommand: "UnexistentCommand",
	}

	evt := &handler.Event{
		K8sEvt: &common.K8sEvent{
			Key:       "key",
			HasSynced: true,
			Kind:      "Add",
			Object:    nil,
		},
		RunNext:     true,
		K8sManifest: []byte("{\"kind\": \"whatever\"}\n"),
		Payload:     []byte{},
	}
	// Adding the event to the storage so the next Update will
	_ = h.Run(context.TODO(), evt)

	// Now we can handler the Update
	evt.K8sEvt.Kind = "Update"

	err := h.Run(context.TODO(), evt)
	if err == nil {
		t.Error("Should have returned an error")
	}
	if evt.RunNext != false {
		t.Error("Should have stopped")
	}

}

func TestDiffTextLines(t *testing.T) {
	diff, err := diffTextLines(
		"diff",
		[]byte("{\"a\": 1}\n"),
		[]byte("{\"a\": 2}\n"),
	)
	if err != nil {
		t.Error(err)
	}

	if len(diff) == 0 {
		t.Errorf("There should have been some difference in the output got %d instead", len(diff))
	}

	diff, err = diffTextLines(
		"diff",
		[]byte("{\"a\": 1}\n"),
		[]byte("{\"a\": 1}\n"),
	)
	if err != nil {
		t.Error(err)
	}

	if len(diff) != 0 {
		t.Errorf("There should have been no differences in the output got %d instead", len(diff))
	}

	diff, err = diffTextLines(
		"unexistentcommand",
		[]byte("{\"a\": 1}\n"),
		[]byte("{\"a\": 1}\n"),
	)
	if err == nil {
		t.Error("Unexistent command should return an error")
	}
}

func TestCreateTempFile(t *testing.T) {
	content := []byte("AnyContent")
	name, err := createTempFile(content)
	if err != nil {
		t.Error(err)
	}
	_content, _ := ioutil.ReadFile(name)
	if !bytes.Equal(content, _content) {
		t.Error("File content mismatch")
	}
}
