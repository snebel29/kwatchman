package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/snebel29/kwatchman/internal/pkg/kwatchman"
	"github.com/snebel29/kwatchman/internal/pkg/cli"
)

func main() {
	c := cli.NewCLI()
	log.Infof("Running kwatchman with %#v", c)
	log.Fatal(kwatchman.Run(c))
}
