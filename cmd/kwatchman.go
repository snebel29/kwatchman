package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/config"
	"github.com/snebel29/kwatchman/internal/pkg/kwatchman"
	"github.com/snebel29/kwatchman/internal/pkg/watcher/k8s"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	w, err := k8s.NewK8sWatcher(conf)
	if err != nil {
		log.Fatal(err)
	}
	if err := kwatchman.Start(w); err != nil {
		log.Fatal(err)
	}
	log.Info("Finishing kwatchman")
}
