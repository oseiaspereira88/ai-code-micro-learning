// Package cli implements the codinho command-line entry points.
package cli

import (
	"context"
	"fmt"
	"io"
)

// Stable process exit codes.
const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

const usage = `codinho - deliberate practice tutor runtime

Usage:
  codinho <command>

Commands:
  version   Print version and build information
  doctor    Diagnose the local installation
  help      Show this help message
`

// Run dispatches args to the matching command and returns the process exit
// code. stdout carries command output; stderr carries usage and errors.
func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	select {
	case <-ctx.Done():
		fmt.Fprintln(stderr, "codinho: cancelled")
		return exitError
	default:
	}

	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return exitUsage
	}

	switch args[0] {
	case "version":
		return runVersion(stdout)
	case "doctor":
		return runDoctor(stdout)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return exitOK
	default:
		fmt.Fprintf(stderr, "codinho: unknown command %q\n\n", args[0])
		fmt.Fprint(stderr, usage)
		return exitUsage
	}
}
