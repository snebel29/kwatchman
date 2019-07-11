package diff

import (
	"encoding/json"
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
