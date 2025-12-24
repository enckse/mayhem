package entities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/enckse/mayhem/internal/state"
)

// LoadJSON will import JSON
func LoadJSON(ctx *state.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	var buf bytes.Buffer
	for scanner.Scan() {
		if _, err := buf.WriteString(scanner.Text()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	var items []Stack
	if err := json.Unmarshal(buf.Bytes(), &items); err != nil {
		return err
	}
	for _, s := range items {
		fmt.Println(s.Title)
		s.Save(ctx)
	}
	return nil
}

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
	s, err := FetchAllStacks(ctx)
	if err != nil {
		return err
	}
	SortStacks(s)
	for _, item := range s {
		tasks := item.Tasks
		SortTasks(tasks)
		item.Tasks = tasks
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	_, err = dst.Write(b)
	return err
}
