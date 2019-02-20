package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/kwatchman"
)

var Version string

func main() {
	log.Fatal(kwatchman.Run(Version))
}
