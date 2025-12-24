package state

import (
	"os"
	"path/filepath"
)

// EnvPrefix is the environment prefix for os.Getenv calls
const EnvPrefix = "MAYHEM_"

func detectDir(xdgName, envVar, altName string) (string, error) {
	p, err := getDir(xdgName, envVar, altName)
	if err != nil {
		return "", err
	}
	return p, os.MkdirAll(p, os.ModePerm)
}

func getDir(xdgName, envVar, altName string) (string, error) {
	path := os.Getenv(EnvPrefix + envVar)
	if path != "" {
		return path, nil
	}

	const appDir = "mayhem"
	xdg := os.Getenv(xdgName)
	if xdg != "" {
		return filepath.Join(xdg, appDir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, altName, appDir), nil
}
