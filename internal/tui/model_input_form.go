package tui

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
	"github.com/enckse/mayhem/internal/tui/keys"
)

type (
	inputForm struct {
		focusIndex    int
		data          entities.Entity
		dataType      string
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

const (
	stackTitleIndex int = iota
)

const (
	taskTitleIndex int = iota
	taskNotesIndex
	taskPriorityIndex
	taskDeadlineIndex
)

const taskLastIndex = taskDeadlineIndex

var (
	stackFields = map[int]field{
		stackTitleIndex: {
			name:             "Title",
			prompt:           "Stack Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         textInputKeys,
			validationPrompt: "Stack title field can not be empty❗",
		},
	}

	taskFields = map[int]field{
		taskTitleIndex: {
			name:             "Title",
			prompt:           "Task Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         textInputKeys,
			validationPrompt: "Task title field can not be empty❗",
		},
		taskNotesIndex: {
			name:     "Notes",
			prompt:   "Task Notes",
			helpKeys: textAreaKeys,
		},
		taskPriorityIndex: {
			name:     "Priority",
			prompt:   "Task Priority",
			helpKeys: listSelectorKeys,
		},
		taskDeadlineIndex: {
			name:     "Deadline",
			prompt:   "Task Deadline",
			helpKeys: timePickerKeys,
		},
	}
)

func initializeInput(selectedTable string, data entities.Entity, fieldIndex int, ctx *state.Context) inputForm {
	var m inputForm
	if selectedTable == "stack" {
		m = inputForm{
			data:       data,
			focusIndex: fieldIndex,
			dataType:   "stack",
			fieldMap:   stackFields,
		}

		targetField := m.fieldMap[fieldIndex]
		stack := data.(entities.Stack)

		switch fieldIndex {
		case 0:
			targetField.model = initializeTextInput(stack.Title, "", 20, goToFormWithVal)
		}

		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField

	} else {
		m = inputForm{
			data:       data,
			focusIndex: fieldIndex,
			fieldMap:   taskFields,
			dataType:   "task",
		}

		targetField := m.fieldMap[fieldIndex]
		task := data.(entities.Task)

		switch fieldIndex {
		case taskTitleIndex:
			targetField.model = initializeTextInput(task.Title, "", 60, goToFormWithVal)
		case taskNotesIndex:
			targetField.model = initializeTextArea(task.Notes, ctx.Screen)
		case taskPriorityIndex:
			opts := []keyVal{
				{val: "0"},
				{val: "1"},
				{val: "2"},
				{val: "3"},
				{val: "4"},
			}
			targetField.model = initializeListSelector(opts, fmt.Sprintf("%d", task.Priority), goToFormWithVal)
		case taskDeadlineIndex:
			if task.Deadline.IsZero() {
				targetField.model = initializeTimePicker(time.Now())
			} else {
				targetField.model = initializeTimePicker(task.Deadline)
			}
		}
		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField
	}

	m.context = ctx
	return m
}

func (m inputForm) Init() tea.Cmd {
	return nil
}

func (m inputForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Transfer control to selectModel's Update method
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Mappings.Return):
			return m, goToMainCmd

		case key.Matches(msg, keys.Mappings.Exit):
			return m, tea.Quit
		}

	case goToFormMsg:
		selectedValue := msg.value

		if (m.fieldMap[m.focusIndex].isRequired) && (selectedValue == m.fieldMap[m.focusIndex].nilValue) {
			m.isInvalid = true
			m.invalidPrompt = m.fieldMap[m.focusIndex].validationPrompt
			return m, nil
		}
		m.isInvalid = false

		switch m.dataType {
		case "stack":
			stack := m.data.(entities.Stack)

			switch m.focusIndex {
			case stackTitleIndex:
				stack.Title = selectedValue.(string)
			}

			stack.Save(m.context)

		case "task":
			task := m.data.(entities.Task)

			switch m.focusIndex {
			case taskTitleIndex:
				task.Title = selectedValue.(string)

				if task.CreatedAt.IsZero() {
					m.isNewTask = true
				}

			case taskNotesIndex:
				task.Notes = selectedValue.(string)
			case taskPriorityIndex:
				task.Priority, _ = strconv.ParseUint(selectedValue.(keyVal).val, 10, 64)
			case taskDeadlineIndex:
				task.Deadline = selectedValue.(time.Time)

			}

			task = task.Save(m.context).(entities.Task)

			if m.isNewTask {
				entities.IncrementPendingCount(task.StackID, m.context)
			}
		}

		return m, goToMainWithVal("refresh")
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
	var cmd tea.Cmd
	inputField := m.fieldMap[m.focusIndex]
	inputField.model, cmd = m.fieldMap[m.focusIndex].model.Update(msg)
	m.fieldMap[m.focusIndex] = inputField

	return m, cmd
}

func (m inputForm) View() string {
	var b strings.Builder

	// ADD changes for invalid input case

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
