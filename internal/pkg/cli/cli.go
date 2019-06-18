package cli

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	Version   string
	namespace = kingpin.Flag(
		"namespace",
		"k8s namespace where to get resources from: default to all").Default(
		"").Envar("KW_NAMESPACE").Short('n').String()
	kubeconfig = kingpin.Flag(
		"kubeconfig",
		"kubeconfig path for running out of k8s").Default(
		fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))).Envar("KW_KUBECONFIG").Short('k').String()
)

// CLIArgs holds the command line arguments
type CLIArgs struct {
	Namespace  string
	Kubeconfig string
}

// NewCLI returns a CLI
func NewCLI() *CLIArgs {
	kingpin.Version(Version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return &CLIArgs{
		Namespace:  *namespace,
		Kubeconfig: *kubeconfig,
	}
}
