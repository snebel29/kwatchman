package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"io/ioutil"
	"runtime"
	"testing"
)

var thisFilename string

func init() {
	_, t, _, _ := runtime.Caller(0)
	thisFilename = t
}

func TestDiffHandlerUseSingletonStorage(t *testing.T) {
	s1 := newStorage()
	addr1 := fmt.Sprintf("%p", s1)
	s2 := newStorage()
	addr2 := fmt.Sprintf("%p", s2)
	if addr1 != addr2 {
		t.Errorf("addr should be the same however %s != %s", addr1, addr2)
	}
}

func TestCleanK8sManifest(t *testing.T) {
	manifest := `
		{
		  "apiVersion": "",
		  "kind": "",
		  "metadata": {
			"generation":2,
			"resourceVersion":"267844584"
		  },
		  "spec": {},
		  "status": {}
		}
	`
	cleaned, _ := cleanK8sManifest([]byte(manifest))
	obj := &k8sObject{}
	json.Unmarshal(cleaned, obj)
	if obj.Metadata.Generation != 0 {
		t.Errorf("Metadata.Generation should be nil, got %v instead", obj.Metadata.Generation)
	}
	if obj.Metadata.ResourceVersion != "" {
		t.Errorf("Metadata.ResourceVersion should be empty, got [%s] instead", obj.Metadata.ResourceVersion)
	}
	if obj.Status != nil {
		t.Errorf("status should be nil, got %#v instead", obj.Status)
	}
}

func TestDiffFuncLogEntryIfthereIsDifferences(t *testing.T) {
	hook := log_test.NewGlobal()

	key := "Key1"
	s := newStorage()

	if _, ok := s[key]; ok {
		t.Errorf("Key %s should NOT exists", key)
	}

	// Fake JSON struct must have some common fields with k8sObject struct
	// In order to unmarshal the differences
	err := DiffFunc(context.TODO(), &common.K8sEvent{
		Key:       key,
		HasSynced: true,
		Kind:      "Update",
		Object:    nil,
	}, []byte("{\"kind\": \"fakeKind\"}\n"))

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	if _, ok := s[key]; !ok {
		t.Errorf("Key %s should exists", key)
	}

	if hook.LastEntry() != nil {
		t.Errorf("Logging lastEntry should be nil")
	}

	err = DiffFunc(context.TODO(), &common.K8sEvent{
		Key:       key,
		HasSynced: true,
		Kind:      "Update",
		Object:    nil,
	}, []byte("{\"kind\": \"fakeKindDifferentThanPrevious\"}\n"))

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	if len(hook.AllEntries()) != 1 || hook.LastEntry().Level.String() != "info" {
		t.Errorf("There should be just one log entry, there is %d instead", len(hook.AllEntries()))
	}
}

func TestDiffTextLines(t *testing.T) {
	diff, err := diffTextLines(
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
		[]byte("{\"a\": 1}\n"),
		[]byte("{\"a\": 1}\n"),
	)
	if err != nil {
		t.Error(err)
	}

	if len(diff) != 0 {
		t.Errorf("There should have been no differences in the output got %d instead", len(diff))
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
