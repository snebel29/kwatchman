package diff

import (
	"context"
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

type diffHandler struct {
	config             config.Handler
	annotationsToClean []string
	storage            *storage
}

// NewDiffHandler return a diff handler and defines the default
// annotations that has to be cleaned to avoid noise due to them chaning on every single event
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

func getObjID(evt *handler.Event) string {
	return fmt.Sprintf("%s/%s", evt.K8sEvt.Key, evt.ResourceKind)
}

// runAdd runs the handler when the event is Add
func (h *diffHandler) runAdd(ctx context.Context, evt *handler.Event) error {

	cleanedManifest := evt.K8sManifest
	h.storage.Add(getObjID(evt), cleanedManifest)

	// Initial sache sync-up events are "Add", we don't want them to be notified
	// but we want them to fill up our storage for future comparison
	if !evt.K8sEvt.HasSynced {
		evt.RunNext = false
	}

	return nil
}

// runUpdate runs the handler when the event is Update
func (h *diffHandler) runUpdate(ctx context.Context, evt *handler.Event) error {

	var diff []byte
	var err error
	cleanedManifest := evt.K8sManifest

	// Since this is an update, there should be a cleaned manifest into the storage
	// for safety we double check, the same apply for HasSynced
	if storedManifest, ok := h.storage.Get(getObjID(evt)); ok && evt.K8sEvt.HasSynced {
		diff, err = diffTextLines(storedManifest, cleanedManifest)
		if err != nil {
			evt.RunNext = false
			return errors.Wrap(err, "diffTextLines")
		}
		// If there is differences we allow for the next handler to run
		// typically a notifier such as slack
		if len(diff) < 0 {
			evt.RunNext = false
		}
		evt.Payload = diff
	}

	// Adding to the storage only after comparison
	h.storage.Add(getObjID(evt), cleanedManifest)
	return nil
}

// runDelete deletes the object from storage and keep moving forward in the chain
func (h *diffHandler) runDelete(ctx context.Context, evt *handler.Event) error {
	h.storage.Delete(getObjID(evt))
	return nil
}

// Run spits out the differentce between previous versions of K8sManifest
// this function is normally the base function handler for resource watchers
// because filters noise by cleaning metadata consolidating logical changes
// from the user perspective, output returns the cleaned manifest and the diff is
// returned in the payload, next handler is run only if a difference is found
func (h *diffHandler) Run(ctx context.Context, evt *handler.Event) error {
	ctx = nil

	switch evt.K8sEvt.Kind {
	case "Add", "Update":
		// Clean only for Add and Update since Delete has no manifest and would fail
		cleanedManifest, err := cleanK8sManifest(evt.K8sManifest, h.annotationsToClean)
		if err != nil {
			evt.RunNext = false
			return err
		}
		evt.K8sManifest = cleanedManifest
	}

	switch evt.K8sEvt.Kind {
	case "Add":
		return h.runAdd(ctx, evt)
	case "Update":
		return h.runUpdate(ctx, evt)
	case "Delete":
		return h.runDelete(ctx, evt)
	}

	// If none of above events matches return an error
	evt.RunNext = false
	return fmt.Errorf("Unknown event kind %s", evt.K8sEvt.Kind)
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
	// which is for now a mandatory requirement, in the future we will be able
	// to configure our own, and eventually make semantic diff using pure go code
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
