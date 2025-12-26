package state_test

import (
	"strings"
	"testing"

	"github.com/enckse/mayhem/internal/state"
)

func TestConfigDatabase(t *testing.T) {
	c := state.Config{}
	c.Data.Directory = "xyz"
	if c.Database() != "xyz/todo.json" {
		t.Errorf("invalid datbase file: %s", c.Database())
	}
}

func TestConfigFile(t *testing.T) {
	cfg := testConfig("settings.toml", t)
	if strings.Contains(cfg.Data.Directory, "~") {
		t.Errorf("invalid data dir: %v", cfg.Data.Directory)
	}
	if cfg.Backups.Directory == "" {
		t.Error("invalid backups dir")
	}
}

func TestConfigEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", ".")
	c := testConfig("", t)
	if c.Data.Directory != "testdata/mayhem" {
		t.Errorf("invalid data dir: %v", c.Data.Directory)
	}
}

func testConfig(file string, t *testing.T) state.Config {
	t.Setenv("XDG_CACHE_HOME", "testdata")
	cfg, err := state.LoadConfig(file)
	if err != nil {
		t.Errorf("invalid load: %v", err)
	}
	return cfg
}

func TestConfigDefaults(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "invalid")
}

func TestConfigBadTOML(t *testing.T) {
	if _, err := state.LoadConfig("config_test.go"); err == nil || err.Error() != "toml: line 1: expected '.' or '=', but got 's' instead" {
		t.Errorf("invalid load: %v", err)
	}
}
