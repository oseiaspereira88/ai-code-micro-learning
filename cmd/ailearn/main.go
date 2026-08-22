// Command ailearn is the administrative CLI entry point for the ailearn
// runtime. It does not implement pedagogical or MCP behavior; see
// internal/cli for the command surface.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/oseiaspereira88/ailearn/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
