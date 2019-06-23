package diff

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/handler"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

var (
	singleton storage
)

type storage map[string][]byte

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
	annotationsToClean []string
}

func NewDiffHandler() handler.Handler {
	return &diffHandler{
		annotationsToClean: []string{
			"deployment.kubernetes.io/revision",
			"kubectl.kubernetes.io/last-applied-configuration",
		},
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

// diff Run spits out the differentce between previous versions of K8sManifest
// this function is normally the base function handler for resource watchers
// because filters noise by cleaning metadata consolidating logical changes
// from the user perspective, output returns the cleaned manifest and the diff is
// returned in the payload
func (h *diffHandler) Run(ctx context.Context, input handler.Input) (handler.Output, error) {
	//TODO: Should cleaning manifest be extracted from DiffFunc into its own Handler?
	ctx = nil
	s := newStorage()

	if input.Evt.Kind == "Delete" {
		delete(s, input.Evt.Key)

		return handler.Output{
			K8sManifest: input.K8sManifest,
			Payload:     input.Payload,
			RunNext:     false, // Delete events won't be handled from now on
		}, nil
	}

	// Only diff if event is Update
	cleanedManifest, err := cleanK8sManifest(input.K8sManifest, h.annotationsToClean)
	if err != nil {
		return handler.Output{
			K8sManifest: input.K8sManifest,
			Payload:     input.Payload,
			RunNext:     false}, err
	}

	var diff []byte
	nextRun := false

	if storedManifest, ok := s[input.Evt.Key]; ok && input.Evt.HasSynced {
		diff, err = diffTextLines(storedManifest, cleanedManifest)
		if err != nil {
			return handler.Output{
				K8sManifest: input.K8sManifest,
				Payload:     input.Payload,
				RunNext:     false}, errors.Wrap(err, "diffTextLines")
		}
		// If there is differences we allow for the next handler to run
		// typically a notifier such as slack
		if len(diff) > 0 {
			nextRun = true
		}
	}

	s[input.Evt.Key] = cleanedManifest
	return handler.Output{
		K8sManifest: cleanedManifest,
		Payload:     diff,
		RunNext:     nextRun,
	}, nil
}

func newStorage() storage {
	lock := &sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()
	if singleton == nil {
		singleton = make(storage)
	}
	return singleton
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
	// TODO: This function is coupled with POSIX diff command
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
