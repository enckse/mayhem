package entities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/enckse/mayhem/internal/state"
)

// LoadJSON will import JSON
func LoadJSON(ctx *state.Context, merge bool, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	var buf bytes.Buffer
	for scanner.Scan() {
		if _, err := buf.WriteString(scanner.Text()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	var existing []Stack
	if merge {
		var err error
		existing, err = FetchStacks(ctx)
		if err != nil {
			return err
		}
	}
	var items []Stack
	if err := json.Unmarshal(buf.Bytes(), &items); err != nil {
		return err
	}
	for _, s := range items {
		if merge {
			if slices.ContainsFunc(existing, func(x Stack) bool {
				return x.Title == s.Title
			}) {
				fmt.Printf("[skipped] %s\n", s.Title)
				continue
			}
		}
		s.Save(ctx)
		fmt.Printf("[imported] %s\n", s.Title)
	}
	return nil
}

// DumpJSONToFile will dump entities to JSON
func DumpJSONToFile(ctx *state.Context) error {
	file, err := os.OpenFile(filepath.Join(ctx.Config.Data.Directory, state.FileName+"json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()
	return DumpJSON(file, ctx)
}

// DumpJSON will write the current JSOn state to stdout
func DumpJSON(dst io.Writer, ctx *state.Context) error {
	s, err := FetchStacks(ctx)
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
