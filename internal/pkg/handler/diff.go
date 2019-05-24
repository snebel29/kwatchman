package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kooper/operator/common"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

var (
	singleton storage
)

type storage map[string]string

// DiffFunc spits out the differentce between two []byte - normally k8s manifests
// and is supposed to be the base function handler for resource watchers
func DiffFunc(_ context.Context, evt *common.K8sEvent, k8sManifest []byte) error {
	s := newStorage()
	if text, ok := s[evt.Key]; ok && evt.HasSynced {

		diff, err := diffTextLines(text, prettyPrintJSON(k8sManifest))
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("%s | %s | %s", evt.Key, evt.Kind, diff)
		}
	}
	s[evt.Key] = prettyPrintJSON(k8sManifest)
	return nil
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

func createTempFile(content string) (string, error) {
	tmpfile, err := ioutil.TempFile("/tmp", "diff-file")
	if err != nil {
		return "", err
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func diffTextLines(text1, text2 string) (string, error) {

	//TODO: Handle concurrency here? could there be race conditions?

	file1, err := createTempFile(text1)
	if err != nil {
		return "", err
	}
	defer os.Remove(file1)

	file2, err := createTempFile(text2)
	if err != nil {
		return "", err
	}
	defer os.Remove(file2)

	output, err := exec.Command("diff", file1, file2).CombinedOutput()
	if err != nil {
		fmt.Println(err, string(output))
		return "", err
	}

	return string(output), nil
}
