package cli

import (
	"fmt"
	"io"
	"runtime/debug"
)

// version and commit are injected at release build time via:
//
//	go build -ldflags "-X .../internal/cli.version=v0.1.0 -X .../internal/cli.commit=<sha>"
//
// Unset, they fall back to the Go module's own VCS stamping so a plain
// `go build` from a Git checkout still reports a real commit without
// requiring Git at runtime.
var (
	version = "dev"
	commit  = "none"
)

func runVersion(stdout io.Writer) int {
	fmt.Fprintf(stdout, "ailearn %s (%s)\n", version, resolveCommit())
	return exitOK
}

func resolveCommit() string {
	if commit != "none" {
		return commit
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return commit
}
