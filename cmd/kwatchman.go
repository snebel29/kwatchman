package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/kwatchman"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
)

func main() {
	log.Fatal(kwatchman.Run(cli.NewCLI()))
}
