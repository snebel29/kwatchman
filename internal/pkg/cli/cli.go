package cli

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Version 	string
	clusterName = kingpin.Flag(
		"cluster-name",
		"Name of k8s cluster where kwatchman is running, use for notification purposes only").Default(
		"undefined").Short('n').String()
)

// CLIArgs holds the command line arguments
type CLIArgs struct {
	ClusterName string
}

// NewCLI returns a CLI
func NewCLI() *CLIArgs {
	kingpin.Version(Version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return &CLIArgs{
		ClusterName: *clusterName,
	}
}