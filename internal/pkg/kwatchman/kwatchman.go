package kwatchman

import (
	"github.com/snebel29/kwatchman/internal/pkg/cli"
	"github.com/snebel29/kwatchman/internal/pkg/watcher"
)

func Run(c *cli.CLIArgs) error {
	w := watcher.NewWatcher(c)
	return w.Run()
}
