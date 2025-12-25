// Package inputs defines forms
package inputs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/inputs/lists"
	"github.com/enckse/mayhem/internal/tui/inputs/text"
	"github.com/enckse/mayhem/internal/tui/inputs/textarea"
	"github.com/enckse/mayhem/internal/tui/inputs/timepicker"
	"github.com/enckse/mayhem/internal/tui/keys"
	"github.com/enckse/mayhem/internal/tui/messages"
)

type (
	// Form is an input form
	Form struct {
		focusIndex    int
		data          entities.Entity
		isStacks      bool
		fieldMap      map[int]field
		isInvalid     bool
		invalidPrompt string
		isNewTask     bool
		helpKeys      keys.Map
		context       *state.Context
	}

	field struct {
		name             string
		prompt           string
		model            tea.Model
		isRequired       bool
		nilValue         string
		validationPrompt string
		helpKeys         keys.Map
	}
)

var (
	stackFields = map[int]field{
		definitions.StackTitleIndex: {
			name:             "Title",
			prompt:           "Stack Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         keys.TextInputMappings,
			validationPrompt: "Stack title field can not be empty❗",
		},
	}

	taskFields = map[int]field{
		definitions.TaskTitleIndex: {
			name:             "Title",
			prompt:           "Task Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         keys.TextInputMappings,
			validationPrompt: "Task title field can not be empty❗",
		},
		definitions.TaskNotesIndex: {
			name:     "Notes",
			prompt:   "Task Notes",
			helpKeys: keys.TextAreaInputMappings,
		},
		definitions.TaskPriorityIndex: {
			name:     "Priority",
			prompt:   "Task Priority",
			helpKeys: keys.ListSelectorMappings,
		},
		definitions.TaskDeadlineIndex: {
			name:     "Deadline",
			prompt:   "Task Deadline",
			helpKeys: keys.TimePickerMappings,
		},
	}
)

// NewStackForm generates a form for stack managing
func NewStackForm(data entities.Entity, ctx *state.Context) Form {
	return newForm(true, data, 0, ctx)
}

// NewTaskForm generates a form for task handling
func NewTaskForm(data entities.Entity, fieldIndex int, ctx *state.Context) Form {
	return newForm(false, data, fieldIndex, ctx)
}

func newForm(isStack bool, data entities.Entity, fieldIndex int, ctx *state.Context) Form {
	var m Form
	if isStack {

		m = Form{
			data:       data,
			focusIndex: fieldIndex,
			fieldMap:   stackFields,
		}

		targetField := m.fieldMap[fieldIndex]
		stack := data.(entities.Stack)

		switch fieldIndex {
		case 0:
			targetField.model = text.New(stack.Title, "", 20, messages.FormGoToWith)
		}

		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField
	} else {
		m = Form{
			data:       data,
			focusIndex: fieldIndex,
			fieldMap:   taskFields,
		}

		targetField := m.fieldMap[fieldIndex]
		task := data.(entities.Task)

		switch fieldIndex {
		case definitions.TaskTitleIndex:
			targetField.model = text.New(task.Title, "", 60, messages.FormGoToWith)
		case definitions.TaskNotesIndex:
			targetField.model = textarea.New(task.Notes, ctx.Screen)
		case definitions.TaskPriorityIndex:
			opts := []definitions.KeyValue{
				{Value: "0"},
				{Value: "1"},
				{Value: "2"},
				{Value: "3"},
				{Value: "4"},
			}
			targetField.model = lists.NewSelector(opts, fmt.Sprintf("%d", task.Priority), messages.FormGoToWith)
		case definitions.TaskDeadlineIndex:
			if task.Deadline.IsZero() {
				targetField.model = timepicker.New(time.Now())
			} else {
				targetField.model = timepicker.New(task.Deadline)
			}
		}
		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField
	}

	m.isStacks = isStack
	m.context = ctx
	return m
}

// HelpKeys will get the help keys for the input form
func (m Form) HelpKeys() keys.Map {
	return m.helpKeys
}

// Init will init the model
func (m Form) Init() tea.Cmd {
	return nil
}

// Update will update the model
func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Transfer control to selectModel's Update method
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Mappings.Return):
			return m, messages.MainGoTo

		case key.Matches(msg, keys.Mappings.Exit):
			return m, tea.Quit
		}

	case messages.Form:
		selectedValue := msg.Value

		if (m.fieldMap[m.focusIndex].isRequired) && (selectedValue == m.fieldMap[m.focusIndex].nilValue) {
			m.isInvalid = true
			m.invalidPrompt = m.fieldMap[m.focusIndex].validationPrompt
			return m, nil
		}
		m.isInvalid = false

		if m.isStacks {
			stack := m.data.(entities.Stack)

			switch m.focusIndex {
			case definitions.StackTitleIndex:
				stack.Title = selectedValue.(string)
			}

			stack.Save(m.context.DB)
		} else {
			task := m.data.(entities.Task)

			switch m.focusIndex {
			case definitions.TaskTitleIndex:
				task.Title = selectedValue.(string)

				if task.CreatedAt.IsZero() {
					m.isNewTask = true
				}

			case definitions.TaskNotesIndex:
				task.Notes = selectedValue.(string)
			case definitions.TaskPriorityIndex:
				task.Priority, _ = strconv.ParseUint(selectedValue.(definitions.KeyValue).Value, 10, 64)
			case definitions.TaskDeadlineIndex:
				task.Deadline = selectedValue.(time.Time)

			}

			task = task.Save(m.context.DB).(entities.Task)
		}

		return m, messages.MainGoToWith("refresh")
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
	var cmd tea.Cmd
	inputField := m.fieldMap[m.focusIndex]
	inputField.model, cmd = m.fieldMap[m.focusIndex].model.Update(msg)
	m.fieldMap[m.focusIndex] = inputField

	return m, cmd
}

// View will display the model
func (m Form) View() string {
	var b strings.Builder

	// NOTE: add changes for invalid input case
	b.WriteString(display.HighlightedTextStyle.Render(m.fieldMap[m.focusIndex].prompt))

	if m.isInvalid {
		b.WriteString(lipgloss.NewStyle().Foreground(display.HighlightedBackgroundColor).Render("    **" + m.invalidPrompt))
	}

	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(m.fieldMap[m.focusIndex].model.View())
	b.WriteRune('\n')
	return b.String()
}
