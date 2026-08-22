package cli

import (
	"fmt"
	"io"
	"runtime"

	"github.com/oseiaspereira88/ailearn/internal/config"
)

// runDoctor reports installation diagnostics only. It never executes files
// from the project workspace.
func runDoctor(stdout io.Writer) int {
	cfg := config.Load()

	fmt.Fprintf(stdout, "go: %s\n", runtime.Version())
	fmt.Fprintf(stdout, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(stdout, "workspace: %s\n", cfg.WorkspaceRoot)
	fmt.Fprintln(stdout, "status: ok")

	return exitOK
}
