// Package convert can export to JSON
package convert

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/enckse/mayhem/internal/app"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

// ToJSON will dump entities to JSON
func ToJSON(ctx *state.Context) error {
	s, err := entities.FetchAllStacks(ctx)
	if err != nil {
		return err
	}
	entities.SortStacks(s)
	for _, item := range s {
		tasks := item.Tasks
		entities.SortTasks(tasks)
		item.Tasks = tasks
	}
	dir, err := app.DataDir()
	if err != nil {
		return err
	}
	payload := filepath.Join(dir, "tasks.json")
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(payload, b, 0o644)
}
