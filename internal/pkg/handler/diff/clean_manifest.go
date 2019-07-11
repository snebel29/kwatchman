package diff

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
)

type k8sObjectMetadata struct {
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`

	Generation int `json:"-"` // generation will be omitted by json.Marshal

	CreationTimestamp string `json:"creationTimestamp"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	ResourceVersion   string `json:"-"` // resourceVersion will be omitted by json.Marshal
	SelfLink          string `json:"selfLink"`
	Uid               string `json:"uid"`
}

type k8sObject struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   k8sObjectMetadata `json:"metadata"`
	Spec       interface{}       `json:"spec"`
	Status     interface{}       `json:"-"` // status will be omitted by json.Marshal
}

func filterMapByKey(m map[string]string, toFilter []string) {
	for _, k := range toFilter {
		delete(m, k)
	}
}

func cleanAnnotations(obj *k8sObject, annotationsToClean []string) {
	if obj.Metadata.Annotations == nil {
		obj.Metadata.Annotations = make(map[string]string)
	}
	filterMapByKey(
		obj.Metadata.Annotations,
		annotationsToClean,
	)
}

// cleanK8sManifest cleans metadata information and indent the manifest in preparation for text
// comparisons
func cleanK8sManifest(manifest []byte, annotationsToClean []string) ([]byte, error) {
	obj := &k8sObject{}

	if err := json.Unmarshal(manifest, obj); err != nil {
		return nil, errors.Wrap(err, "cleanK8sManifest Unmarshal")
	}

	cleanAnnotations(obj, annotationsToClean)

	_cleanK8sManifest, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "cleanK8sManifest Marshal")
	}

	_json, err := handler.PrettyPrintJSON(_cleanK8sManifest)
	if err != nil {
		return nil, errors.Wrap(err, "cleanK8sManifest prettyPrintJSON")
	}

	return _json, nil
}
