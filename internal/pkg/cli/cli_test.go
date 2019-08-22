package cli

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCLIWrongArgsExitCode(t *testing.T) {
	expectedExitCode := 1
	if os.Getenv("TESTING") == "true" {
		NewCLI()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCLIWrongArgsExitCode")
	cmd.Env = append(os.Environ(), "TESTING=true")
	err := cmd.Run()
	if err != nil && err.Error() == fmt.Sprintf("exit status %d", expectedExitCode) {
		return
	}
	t.Errorf("process ran with %v, want exit status %d", err, expectedExitCode)
}

func TestCliArgs(t *testing.T) {
	namespace := "myNamespace"
	kubeconfig := "myKubeconfig"
	config := "myConfig"
	labels := "environment,environment notin (frontend)"
	logLevel := "info"
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--namespace=%s", namespace),
		fmt.Sprintf("--kubeconfig=%s", kubeconfig),
		fmt.Sprintf("--config=%s", config),
		fmt.Sprintf("--label-selector=%s", labels),
		fmt.Sprintf("--log-level=%s", logLevel),
	}

	cli := NewCLI()
	if cli.Namespace != namespace {
		t.Errorf("%s != %s", cli.Namespace, namespace)
	}
	if cli.Kubeconfig != kubeconfig {
		t.Errorf("%s != %s", cli.Kubeconfig, kubeconfig)
	}
	if cli.ConfigFile != config {
		t.Errorf("%s != %s", cli.ConfigFile, config)
	}
	if cli.LabelSelector != labels {
		t.Errorf("%s != %s", cli.LabelSelector, labels)
	}
	if cli.LogLevel != logLevel {
		t.Errorf("%s != %s", cli.LogLevel, logLevel)
	}
}
