package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/tui/inputs/timepicker"
	"github.com/enckse/mayhem/internal/tui/keys"
)

var (
	stackKeys = keys.Map{
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("'n'", "new stack 🌟"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit 📝"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("'x'", "delete 🗑"),
		),
	}

	taskKeys = keys.Map{
		Toggle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("'tab'", "check/uncheck 🔄"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("'n'", "new task 🌟"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit 📝"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("'x'", "delete 🗑"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("'m'", "change stack 📤"),
		),
	}

	tableNavigationKeys = keys.Map{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("'↑/k'", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("'↓/j'", "down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("'g'", "jump to top"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("'G'", "jump to bottom"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("'←/h'", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("'→/l'", "right"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("'?'", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("'q'", "quit"),
		),
	}

	taskFinishStatus = map[uint]bool{}
)

func stackColumns() []table.Column {
	return []table.Column{
		{Title: "       Stacks 🗃", Width: 20},
		{Title: "", Width: 2},
	}
}

func taskColumns() []table.Column {
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "           Tasks 📄", Width: 30},
		{Title: "     Deadline 🕑", Width: 20},
		{Title: "Priority", Width: 8},
	}
}

func stackRows(stacks []entities.Stack) []table.Row {
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

func taskRows(tasks []entities.Task) []table.Row {
	rows := make([]table.Row, len(tasks))

	// We perform this step earlier since we need the deadline & finish status data before sorting
	for _, val := range tasks {
		taskFinishStatus[val.ID] = val.IsFinished
	}

	entities.SortTasks(tasks)

	var prefix string
	var deadline string

	for i, val := range tasks {
		deadline = formatTime(val.Deadline, true)

		if taskFinishStatus[val.ID] {
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

func buildTable(columns []table.Column, tableType display.TableType, screen *display.Screen) table.Model {
	t := table.New(
		table.WithHeight(screen.Table.ViewHeight),
		table.WithColumns(columns),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	s := display.TableStyle(tableType)
	t.SetStyles(s)

	return t
}

func formatTime(time time.Time, fullDate bool) string {
	if time.IsZero() {
		return fmt.Sprintf("%10s", "-")
	}

	year := fmt.Sprintf("%04d", time.Year())
	month := fmt.Sprintf("%02d", int(time.Month()))
	days := fmt.Sprintf("%02d", time.Day())
	hours := fmt.Sprintf("%02d", timepicker.FormatHour(time.Hour()))
	minutes := fmt.Sprintf("%02d", time.Minute())
	midDayInfo := timepicker.RenderMidDayInfo(time.Hour())

	if fullDate {
		return days + "-" + month + "-" + year + "  " + hours + ":" + minutes + " " + midDayInfo
	}
	return hours + ":" + minutes + " " + midDayInfo
}

func incompleteTaskTag(count uint64) string {
	if count > 0 && count <= 10 {
		return " " + string(rune('➊'+count-1))
	} else if count > 10 {
		return "+➓"
	}
	return ""
}

func findIndex[T interface{ EntityID() uint }](arr []T, id uint) int {
	for i, val := range arr {
		if val.EntityID() == id {
			return i
		}
	}
	return -1
}
