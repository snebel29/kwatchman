package cli

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	Version     string
	clusterName = kingpin.Flag(
		"cluster-name",
		"Name of k8s cluster where kwatchman is running, use for notification purposes only").Default(
		"undefined").Short('c').String()
	namespace = kingpin.Flag(
		"namespace",
		"k8s namespace where to get resources from: default to all").Default(
		"").Short('n').String()
	kubeconfig = kingpin.Flag(
		"kubeconfig",
		"kubeconfig path for running out of k8s").Default(
		fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))).Short('k').String()
)

// CLIArgs holds the command line arguments
type CLIArgs struct {
	ClusterName string
	Namespace   string
	Kubeconfig  string
}

// NewCLI returns a CLI
func NewCLI() *CLIArgs {
	kingpin.Version(Version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return &CLIArgs{
		ClusterName: *clusterName,
		Namespace:   *namespace,
		Kubeconfig:  *kubeconfig,
	}
}
