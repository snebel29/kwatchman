package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/snebel29/kwatchman/internal/pkg/kwatchman"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

func main() {
	c := cli.NewCLI()
	log.Infof("Running kwatchman with %#v", c)
	w, err := watcher.NewK8sWatcher(c)
	if err != nil {
		log.Fatal(err)
	}
	if err := kwatchman.Start(w); err != nil {
		log.Fatal(err)
	}
	log.Info("Finishing kwatchman")
}
