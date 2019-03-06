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
	if err.Error() == fmt.Sprintf("exit status %d", expectedExitCode) {
		return
	}
	t.Errorf("process ran with %v, want exit status %d", err, expectedExitCode)
}

func TestCliArgs(t *testing.T) {
	clusterName := "myClusterName"
	os.Args = []string{
		"kwatchman",
		fmt.Sprintf("--cluster-name=%s", clusterName),
	}

	cli := NewCLI()
	if cli.ClusterName != clusterName {
		t.Errorf("%s != %s", cli.ClusterName, clusterName)
	}
}
