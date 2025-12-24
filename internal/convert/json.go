// Package convert can export to JSON
package convert

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

// ToJSON will dump entities to JSON
func ToJSON(ctx *state.Context) error {
	file, err := os.OpenFile(filepath.Join(ctx.Config.Data.Directory, "tasks.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	return handleJSON(file, ctx)
}

// DumpJSON will write the current JSOn state to stdout
func DumpJSON(ctx *state.Context) error {
	return handleJSON(os.Stdout, ctx)
}

func handleJSON(dst io.Writer, ctx *state.Context) error {
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
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	_, err = dst.Write(b)
	return err
}
