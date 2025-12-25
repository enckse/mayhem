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
	if err := cfg.Backup(time.Now(), time.Now()); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	os.MkdirAll(cfg.Data.Directory, 0o755)
	os.WriteFile(cfg.Database(), []byte{}, 0o644)
	var zero time.Time
	if err := cfg.Backup(now, zero); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ := os.ReadDir(cfg.Backups.Directory)
	if len(children) != 2 {
		t.Errorf("invalid children: %v", children)
	}
	if err := cfg.Backup(now.Add(1*time.Second), time.Now().Add(1*time.Second)); err != nil {
		t.Errorf("unexpected backup error: %v", err)
	}
	children, _ = os.ReadDir(cfg.Backups.Directory)
	if len(children) != 2 {
		t.Errorf("invalid children: %v", children)
	}
}
