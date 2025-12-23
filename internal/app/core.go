// Package app provides common application definitions
package app

import (
	"os"
	"path/filepath"
)

// EnvPrefix is the environment prefix for os.Getenv calls
const EnvPrefix = "MAYHEM_"

// DataDir will get the data directory for db storage
func DataDir() (string, error) {
	p, err := detectDir()
	if err != nil {
		return "", err
	}
	return p, os.MkdirAll(p, os.ModePerm)
}

func detectDir() (string, error) {
	path := os.Getenv(EnvPrefix + "DB_PATH")
	if path != "" {
		return path, nil
	}

	const appDir = "mayhem"
	xdg := os.Getenv("XDG_CACHE_HOME")
	if xdg != "" {
		return filepath.Join(xdg, appDir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", appDir), nil
}
