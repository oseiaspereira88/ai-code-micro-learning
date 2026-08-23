// Package config resolves runtime defaults for the codinho CLI without
// reading secrets, environment files, or the network.
package config

import "os"

// Config holds resolved runtime defaults.
type Config struct {
	WorkspaceRoot string
}

// Load resolves Config from process defaults only.
func Load() Config {
	root, err := os.Getwd()
	if err != nil {
		root = "."
	}
	return Config{WorkspaceRoot: root}
}
