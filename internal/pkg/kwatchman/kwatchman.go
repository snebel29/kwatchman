package kwatchman

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
	"os"
	"os/signal"
	"syscall"
)

func Run(w watcher.Watcher) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-shutdown
		log.Infof("Shutdown upon signal %s", sig.Signal)
		w.Shutdown()
		os.Exit(0)
	}()

	return w.Run()
}
