package diff

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"github.com/snebel29/kwatchman/internal/pkg/registry"
	"io/ioutil"
	"os"
	"os/exec"
)

func init() {
	registry.Register(registry.HANDLER, "diff", NewDiffHandler)
}

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

type diffHandler struct {
	config             config.Handler
	annotationsToClean []string
	storage            *storage
}

func NewDiffHandler(c config.Handler) handler.Handler {
	return &diffHandler{
		config: c,
		annotationsToClean: []string{
			"deployment.kubernetes.io/revision",
			"kubectl.kubernetes.io/last-applied-configuration",
		},
		storage: newStorage(),
	}
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

func stopHandlerWithError(input handler.Input, err error) (handler.Output, error) {
	return handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     false}, err
}

func getObjID(input handler.Input) string {
	return fmt.Sprintf("%s/%s", input.Evt.Key, input.ResourceKind)
}

func (h *diffHandler) runAdd(ctx context.Context, input handler.Input) (handler.Output, error) {

	cleanedManifest, err := cleanK8sManifest(input.K8sManifest, h.annotationsToClean)
	if err != nil {
		return stopHandlerWithError(input, err)
	}

	h.storage.Add(getObjID(input), cleanedManifest)

	runNext := true
	// Initial scache ync-up events are "Add", we don't want them to be notified
	// but we want them to fill up our storage for future comparison
	if !input.Evt.HasSynced {
		runNext = false
	}

	return handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     runNext,
	}, nil
}

func (h *diffHandler) runUpdate(ctx context.Context, input handler.Input) (handler.Output, error) {

	cleanedManifest, err := cleanK8sManifest(input.K8sManifest, h.annotationsToClean)
	if err != nil {
		return stopHandlerWithError(input, err)
	}

	var diff []byte
	runNext := false

	// Since this is an update, there should be a cleaned manifest into the storage
	// for safety we double check, the same apply for HasSynced
	if storedManifest, ok := h.storage.Get(getObjID(input)); ok && input.Evt.HasSynced {
		diff, err = diffTextLines(storedManifest, cleanedManifest)
		if err != nil {
			return stopHandlerWithError(input, errors.Wrap(err, "diffTextLines"))
		}
		// If there is differences we allow for the next handler to run
		// typically a notifier such as slack
		if len(diff) > 0 {
			runNext = true
		}
	}

	// Adding to the storage only after comparison
	h.storage.Add(getObjID(input), cleanedManifest)

	return handler.Output{
		K8sManifest: cleanedManifest,
		Payload:     diff,
		RunNext:     runNext,
	}, nil
}

// runDelete deletes the object from storage and keep moving forward in the chain
func (h *diffHandler) runDelete(ctx context.Context, input handler.Input) (handler.Output, error) {

	h.storage.Delete(getObjID(input))
	return handler.Output{
		K8sManifest: input.K8sManifest,
		Payload:     input.Payload,
		RunNext:     true,
	}, nil
}

// Run spits out the differentce between previous versions of K8sManifest
// this function is normally the base function handler for resource watchers
// because filters noise by cleaning metadata consolidating logical changes
// from the user perspective, output returns the cleaned manifest and the diff is
// returned in the payload, next handler is run only if a difference is found
func (h *diffHandler) Run(ctx context.Context, input handler.Input) (handler.Output, error) {
	ctx = nil

	switch input.Evt.Kind {
	case "Add":
		return h.runAdd(ctx, input)
	case "Update":
		return h.runUpdate(ctx, input)
	case "Delete":
		return h.runDelete(ctx, input)
	}

	return stopHandlerWithError(input, fmt.Errorf("Unknown event kind %s", input.Evt.Kind))
}

func createTempFile(content []byte) (string, error) {
	tmpfile, err := ioutil.TempFile("/tmp", "diff-file")
	if err != nil {
		return "", err
	}
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func diffTextLines(text1, text2 []byte) ([]byte, error) {
	// This function is currently coupled with POSIX diff command
	// which is a mandatory requirement
	file1, err := createTempFile(text1)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file1)

	file2, err := createTempFile(text2)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file2)

	output, err := exec.Command("diff", file1, file2).CombinedOutput()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// Do nothing, this is expected to have exit code

		default:
			log.Errorf("%s %s", err, string(output))
			return nil, err
		}
	}

	return output, nil
}
