package diff

import (
	"bytes"
	"context"
	"encoding/json"
	log_test "github.com/sirupsen/logrus/hooks/test"
	"github.com/snebel29/kooper/operator/common"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"io/ioutil"
	"reflect"
	"testing"
)

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
	annotationsToClean := []string{
		"deployment.kubernetes.io/revision",
		"kubectl.kubernetes.io/last-applied-configuration",
	}
	cleaned, _ := cleanK8sManifest([]byte(manifest), annotationsToClean)
	obj := &k8sObject{}
	err := json.Unmarshal(cleaned, obj)
	if err != nil {
		t.Error(err)
	}
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

func TestDiffHandler(t *testing.T) {
	hook := log_test.NewGlobal()

	// Fake JSON struct must have some common fields with k8sObject struct
	// In order to unmarshal the differences

	// A new key is added, no difference, no error, and nextRun == true should be expected
	h1 := NewDiffHandler(config.Handler{})
	key := "key1"

	output, err := h1.Run(
		context.TODO(),
		handler.Input{
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
	if len(diff) != 0 {
		t.Error("There should be no diff because is a new event")
	}

	if output.RunNext != false {
		t.Error("nextRun should be false because there is no difference")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	if hook.LastEntry() != nil {
		t.Errorf("Logging lastEntry should be nil")
	}

	// The same key with diffrent kind is Updated on the same handler
	// now a difference should be returned
	output, err = h1.Run(
		context.TODO(),
		handler.Input{
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

	if output.RunNext != true {
		t.Error("nextRun should be true because difference was found")
	}

	if err != nil {
		t.Error("No error should have ocurred on Diff")
	}

	// Manifest should be deleted from storage and next handlers not to be trigger
	output, err = h1.Run(
		context.TODO(),
		handler.Input{
			Evt: &common.K8sEvent{
				Key:       key,
				HasSynced: true,
				Kind:      "Delete",
				Object:    nil,
			},
			K8sManifest: []byte{},
			Payload:     []byte{},
		},
	)

	diff = output.Payload
	if len(diff) != 0 {
		t.Error("There should be no difference")
	}

	if output.RunNext != false {
		t.Error("nextRun should be false because there is no difference")
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
	annotationsToClean := []string{
		"deployment.kubernetes.io/revision",
		"kubectl.kubernetes.io/last-applied-configuration",
	}

	for _, annotation := range annotationsToClean {
		m1[annotation] = "whatever"
	}

	obj.Metadata.Annotations = m1
	cleanAnnotations(obj, annotationsToClean)

	m2 := map[string]string{}

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("k8sObject Annotations should be clean!, got %#v instead", m1)
	}
}

func TestCleanAnnotationsWorksWithNONInitializedAnnotations(t *testing.T) {
	obj := &k8sObject{}
	annotationsToClean := []string{
		"deployment.kubernetes.io/revision",
		"kubectl.kubernetes.io/last-applied-configuration",
	}

	cleanAnnotations(obj, annotationsToClean)
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
