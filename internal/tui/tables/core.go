// Package tables provides help to create UI table outputs
package tables

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/inputs/timepicker"
)

var (
	// StackColumns are the table columns for stacks
	StackColumns = []table.Column{
		{Title: "       Stacks 🗃", Width: 20},
		{Title: "", Width: 2},
	}
	// TaskColumns are the table columns for tasks
	TaskColumns = []table.Column{
		{Title: "", Width: 1},
		{Title: "           Tasks 📄", Width: 30},
		{Title: "     Deadline 🕑", Width: 20},
		{Title: "Priority", Width: 8},
	}
)

// StackRows will generate rows for stack
func StackRows(stacks []entities.Stack) []table.Row {
	rows := make([]table.Row, len(stacks))

	entities.SortStacks(stacks)

	for i, val := range stacks {
		row := []string{
			val.Title,
			incompleteTaskTag(val.PendingTaskCount),
		}
		rows[i] = row
	}
	return rows
}

// TaskRows will generate rows for tasks
func TaskRows(tasks []entities.Task, ctx *state.Context) []table.Row {
	rows := make([]table.Row, len(tasks))

	// We perform this step earlier since we need the deadline & finish status data before sorting
	for _, val := range tasks {
		ctx.State.FinishedTasks[val.ID] = val.IsFinished
	}

	entities.SortTasks(tasks)

	var prefix string
	var deadline string

	for i, val := range tasks {
		deadline = timepicker.FormatTime(val.Deadline, true)

		if ctx.State.FinishedTasks[val.ID] {
			prefix = "✘"
		} else {
			prefix = "▢"
		}

		row := []string{
			prefix,
			val.Title,
			deadline,
			fmt.Sprintf("   %d", val.Priority),
		}

		rows[i] = row
	}

	return rows
}

// New will generate a new table model
func New(columns []table.Column, tableType display.TableType, screen *display.Screen) table.Model {
	t := table.New(
		table.WithHeight(screen.Table.ViewHeight),
		table.WithColumns(columns),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	s := display.TableStyle(tableType)
	t.SetStyles(s)

	return t
}

func incompleteTaskTag(count uint64) string {
	if count > 0 && count <= 10 {
		return " " + string(rune('➊'+count-1))
	} else if count > 10 {
		return "+➓"
	}
	return ""
}
