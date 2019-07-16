package cli

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	// Version is initialized using LD_FLAGS during build process (Makefile)
	Version   string
	namespace = kingpin.Flag(
		"namespace",
		"k8s namespace where to get resources from: default to all").Default(
		"").Envar("KW_NAMESPACE").Short('n').String()
	kubeconfig = kingpin.Flag(
		"kubeconfig",
		"kubeconfig path for running out of k8s").Default(
		fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))).Envar("KW_KUBECONFIG").Short('k').String()
	configFile = kingpin.Flag(
		"config",
		"The kwatchman config file").Default(
		fmt.Sprintf("%s/config.toml", os.Getenv("PWD"))).Envar("KW_CONFIG_FILE").Short('c').String()
	labelSelector = kingpin.Flag(
		"label-selector",
		"k8s label selector string: default to all").Default(
		"").Envar("KW_LABEL_SELECTOR").Short('l').String()
)

// Args holds the command line arguments
type Args struct {
	Namespace     string
	Kubeconfig    string
	ConfigFile    string
	LabelSelector string
}

// NewCLI returns a CLI
func NewCLI() *Args {
	kingpin.Version(Version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return &Args{
		Namespace:     *namespace,
		Kubeconfig:    *kubeconfig,
		ConfigFile:    *configFile,
		LabelSelector: *labelSelector,
	}
}
