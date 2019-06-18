package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"io/ioutil"
	"reflect"
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
			"resourceVersion":"267844584",
			"annotations":{"deployment.kubernetes.io/revision": "1"}
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
	if !reflect.DeepEqual(obj.Metadata.Annotations, map[string]string{}) {
		t.Errorf("Metadata.Annotations should be nil, got %#v instead", obj.Metadata.Annotations)
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
	output, err := DiffFunc(
		context.TODO(),
		Input{
			Evt: &common.K8sEvent{
				Key:       key,
				HasSynced: true,
				Kind:      "Add",
				Object:    nil,
			},
			K8sManifest: []byte("{\"kind\": \"fakeKind\"}\n"),
			Payload:     []byte{},
		},
	)

	diff := output.Payload
	if reflect.DeepEqual(diff, []byte{}) {
		t.Error("diff should be empty")
	}

	if output.RunNext != false {
		t.Error("nextRun should be false")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	if _, ok := s[key]; !ok {
		t.Errorf("Key %s should exists", key)
	}

	if hook.LastEntry() != nil {
		t.Errorf("Logging lastEntry should be nil")
	}

	// In this case a difference should be raised
	output, err = DiffFunc(
		context.TODO(),
		Input{
			Evt: &common.K8sEvent{
				Key:       key,
				HasSynced: true,
				Kind:      "Update",
				Object:    nil,
			},
			K8sManifest: []byte("{\"kind\": \"fakeKindDifferentThanPrevious\"}\n"),
			Payload:     []byte{},
		},
	)

	diff = output.Payload
	if len(diff) < 1 {
		t.Error("there should be some difference")
	}

	// Because differences trigger next handler
	if output.RunNext != true {
		t.Error("nextRun should be true")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
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

func TestCleanAnnotationsWorksWithInitializedAnnotations(t *testing.T) {
	obj := &k8sObject{}
	m1 := make(map[string]string)

	for _, annotation := range AnnotationsToClean {
		m1[annotation] = "whatever"
	}

	obj.Metadata.Annotations = m1
	cleanAnnotations(obj)

	m2 := map[string]string{}

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("k8sObject Annotations should be clean!, got %#v instead", m1)
	}
}

func TestCleanAnnotationsWorksWithNONInitializedAnnotations(t *testing.T) {
	obj := &k8sObject{}

	cleanAnnotations(obj)
	m2 := map[string]string{}

	if !reflect.DeepEqual(obj.Metadata.Annotations, m2) {
		t.Errorf(
			"k8sObject Annotations should have been automatically initialized after cleaning it!, got %#v instead", obj.Metadata.Annotations)
	}
}

func TestFilterMapByKey(t *testing.T) {
	m1 := map[string]string{"a": "1", "b": "2", "c": "3"}
	toFilter := []string{"a", "c", "z"}
	filterMapByKey(m1, toFilter)
	m2 := map[string]string{"b": "2"}
	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Maps should match, %#v != %#v", m1, m2)
	}
}
