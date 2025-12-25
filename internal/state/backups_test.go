package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/state"
)

func TestBackup(t *testing.T) {
	now := time.Now()
	cfg := state.Config{}
	cfg.Backups.Directory = filepath.Join("testdata", "backups")
	cfg.Data.Directory = cfg.Backups.Directory
	os.RemoveAll(cfg.Backups.Directory)
	cfg.Backups.Duration = "1h"
	if err := cfg.Backup(time.Now()); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	os.MkdirAll(cfg.Data.Directory, 0o755)
	os.WriteFile(cfg.Database(), []byte{}, 0o644)
	cfg.Backups.Duration = ""
	if err := cfg.Backup(now); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ := os.ReadDir(cfg.Backups.Directory)
	if len(children) != 2 {
		t.Errorf("invalid children: %v", children)
	}
	cfg.Backups.Duration = "1s"
	if err := cfg.Backup(now.Add(1 * time.Second)); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ = os.ReadDir(cfg.Backups.Directory)
	if len(children) != 2 {
		t.Errorf("invalid children: %v", children)
	}
	cfg.Backups.Format = "2006"
	cfg.Backups.Duration = "5h"
	if err := cfg.Backup(now.Add(1 * time.Second)); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ = os.ReadDir(cfg.Backups.Directory)
	if len(children) != 3 {
		t.Errorf("invalid children: %v", children)
	}
	if err := cfg.Backup(now.Add(1 * time.Second)); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ = os.ReadDir(cfg.Backups.Directory)
	if len(children) != 3 {
		t.Errorf("invalid children: %v", children)
	}
}
