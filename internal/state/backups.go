package state

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Backup will perform a backup based on the configuration
func (c Config) Backup(timestamp, threshold time.Time) error {
	db := c.Database()
	if !PathExists(db) {
		return nil
	}
	if err := os.MkdirAll(c.Backups.Directory, os.ModePerm); err != nil {
		return err
	}
	if !threshold.IsZero() {
		err := filepath.Walk(c.Backups.Directory, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(p, "."+databaseName) {
				if info.ModTime().Before(threshold) {
					if err := os.Remove(p); err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	target := timestamp.Format("20060102T150405")
	target = filepath.Join(c.Backups.Directory, fmt.Sprintf("%s.%s", target, databaseName))
	if PathExists(target) {
		return nil
	}
	if err := exec.Command("cp", c.Database(), target).Run(); err != nil {
		return err
	}
	return nil
}
