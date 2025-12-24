// Package ui handles the TUI
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/deletion"
	"github.com/enckse/mayhem/internal/tui/details"
	"github.com/enckse/mayhem/internal/tui/help"
	"github.com/enckse/mayhem/internal/tui/inputs"
	"github.com/enckse/mayhem/internal/tui/inputs/lists"
	"github.com/enckse/mayhem/internal/tui/keys"
	"github.com/enckse/mayhem/internal/tui/messages"
	"github.com/enckse/mayhem/internal/tui/tables"
)

type (
	// ModelWrapper provides a wrapper around a backing TUI model
	ModelWrapper struct {
		Backing *model
	}

	model struct {
		data            []entities.Stack
		stackTable      table.Model
		taskTable       table.Model
		taskDetails     details.Box
		help            help.Model
		input           inputs.Form
		showTasks       bool
		showDetails     bool
		showInput       bool
		showHelp        bool
		customInput     tea.Model
		customInputType string
		showCustomInput bool
		navigationKeys  keys.Map
		preInputFocus   string // useful for reverting back when input box is closed
		firstRender     bool
		prevState       preserveState
		context         *state.Context
	}

	preserveState struct {
		retainState bool
		stackID     uint
		taskID      uint
	}
)

const (
	stackViewName     = "stack"
	detailViewName    = "detail"
	taskViewName      = "task"
	tasksDataUpdate   = "tasks"
	stacksDataUpdate  = "stacks"
	detailsDataUpdate = "details"
)

// Initialize will startup the core application model
func Initialize(ctx *state.Context) ModelWrapper {
	stacks, _ := entities.FetchStacks(ctx)

	m := &model{
		stackTable: tables.New(tables.StackColumns, display.StackTableType, ctx.Screen),
		taskTable:  tables.New(tables.TaskColumns, display.TaskTableType, ctx.Screen),
		// we can't build the details box at this stage since we need both stack & task indices for that
		taskDetails:    details.NewBox(ctx.Screen),
		data:           stacks,
		help:           help.NewModel(keys.StackMappings),
		navigationKeys: keys.TableMappings,
		showHelp:       true,
		context:        ctx,
	}

	m.stackTable.Focus()
	m.taskTable.Blur()
	m.taskDetails.Blur()
	return ModelWrapper{m}
}

// Init initializes the model
func (m *model) Init() tea.Cmd {
	m.firstRender = true
	return nil
}

// Update will update the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Transfer control to inputForm's Update method
	if m.showInput {
		switch msg := msg.(type) {

		case messages.Main:
			m.input = inputs.Form{}
			m.showInput = false

			if msg.Value.(string) == "refresh" {
				m.preserveState()
				m.refreshData()
			}

			switch m.preInputFocus {
			case stackViewName, taskViewName:
				if m.preInputFocus == stackViewName {

					m.stackTable.Focus()
					m.help = help.NewModel(keys.StackMappings)
				} else {
					m.taskTable.Focus()
					m.help = help.NewModel(keys.TaskMappings)
				}
				m.navigationKeys = keys.TableMappings
			case detailViewName:
				m.taskDetails.Focus()
				m.navigationKeys = keys.DetailsMappings
			}

			m.updateViewDimensions(10)

			return m, nil

		case tea.WindowSizeMsg:
			m.context.Screen.Width = msg.Width
			m.context.Screen.Height = msg.Height
			m.updateViewDimensions(14)
			return m, nil

		default:
			inp, cmd := m.input.Update(msg)
			t, _ := inp.(inputs.Form)
			m.input = t

			return m, cmd
		}
	}

	if m.showCustomInput {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.context.Screen.Width = msg.Width
			m.context.Screen.Height = msg.Height
			m.updateViewDimensions(14)
			return m, nil
		}

		switch m.customInputType {
		// Transfer control to delete confirmation model
		case "delete":
			switch msg := msg.(type) {

			case messages.Main:
				m.showCustomInput = false

				switch m.preInputFocus {
				case stackViewName:
					m.stackTable.Focus()
					m.help = help.NewModel(keys.StackMappings)
				case taskViewName:
					m.taskTable.Focus()
					m.help = help.NewModel(keys.TaskMappings)
				}

				if msg.Value.(string) == "y" {
					switch m.preInputFocus {
					case stackViewName:
						stackIndex := m.stackTable.Cursor()
						currStack := m.data[stackIndex]

						if stackIndex == len(m.stackTable.Rows())-1 {
							m.stackTable.SetCursor(stackIndex - 1)
						}

						currStack.Delete(m.context)
						m.showTasks = false
						m.showDetails = false
						m.refreshData()
						return m, nil

					case taskViewName:
						stackIndex := m.stackTable.Cursor()
						taskIndex := m.taskTable.Cursor()

						var currTask entities.Task
						if len(m.data[stackIndex].Tasks) > 0 {
							currTask = m.data[stackIndex].Tasks[taskIndex]

							if !currTask.IsFinished {
								stack := m.data[stackIndex]
								stack.PendingTaskCount--
								stack.Save(m.context)
							}
							if taskIndex == len(m.taskTable.Rows())-1 {
								m.taskTable.SetCursor(taskIndex - 1)
							}
							currTask.Delete(m.context)
							m.refreshData()
							return m, nil
						}
					}
				}

			default:
				inp, cmd := m.customInput.Update(msg)
				t, _ := inp.(deletion.Confirmation)
				m.customInput = t

				return m, cmd
			}

		case "move":
			switch msg := msg.(type) {

			case messages.Main:
				m.showCustomInput = false
				m.taskTable.Focus()
				m.help = help.NewModel(keys.TaskMappings)

				response := msg.Value.(definitions.KeyValue)

				if response.Value == "" {
					return m, nil
				}

				newStackID := response.Key

				stackIndex := m.stackTable.Cursor()
				taskIndex := m.taskTable.Cursor()

				currStack := m.data[stackIndex]
				currTask := currStack.Tasks[taskIndex]

				if currTask.StackID == newStackID {
					return m, nil
				}

				// Moving recurring tasks wouldn't have any effect on the stack pending task count

				// Decrease pending task count for old stack
				if !currTask.IsFinished {
					currStack.PendingTaskCount--
					currStack.Save(m.context)
				}

				// Increase pending task count for new stack
				entities.IncrementPendingCount(newStackID, m.context)
				currTask.StackID = newStackID
				currTask.Save(m.context)

				if taskIndex == len(m.taskTable.Rows())-1 {
					m.taskTable.SetCursor(taskIndex - 1)
				}
				m.refreshData()
				return m, nil

			default:
				inp, cmd := m.customInput.Update(msg)
				t, _ := inp.(lists.Selector)
				m.customInput = t

				return m, cmd
			}
		}
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		// Inter-table navigation
		case key.Matches(msg, keys.Mappings.Left):
			if m.stackTable.Focused() {
				if m.showDetails {
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = help.NewModel(keys.TaskDetailsMappings)
					m.navigationKeys = keys.DetailsMappings

				}
			} else if m.taskTable.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = help.NewModel(keys.StackMappings)
				m.navigationKeys = keys.TableMappings

			} else if m.taskDetails.Focused() {
				m.stackTable.Blur()
				m.taskTable.Focus()
				m.taskDetails.Blur()
				m.help = help.NewModel(keys.TaskMappings)
				m.navigationKeys = keys.TableMappings

			}
			return m, nil

		case key.Matches(msg, keys.Mappings.Right):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) > 0 {
					m.showTasks = true
					m.stackTable.Blur()
					m.taskTable.Focus()
					m.taskDetails.Blur()
					m.help = help.NewModel(keys.TaskMappings)
					m.navigationKeys = keys.TableMappings
					return m, nil
				}
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = help.NewModel(keys.TaskDetailsMappings)
					m.navigationKeys = keys.DetailsMappings
					return m, nil
				}
			} else if m.taskDetails.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = help.NewModel(keys.StackMappings)
				m.navigationKeys = keys.TableMappings
				return m, nil
			}

		// Intra-table navigation

		// When we switch to a new stack:
		//		- Empty task box is shown
		//		- Details box is hidden

		// When we switch to a new task:
		//		- Empty details box is shown
		case key.Matches(msg, keys.Mappings.Up):
			if m.stackTable.Focused() {
				m.stackTable.MoveUp(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.FocusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData(tasksDataUpdate)
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveUp(1)
				m.taskDetails.FocusIndex = 0
				m.showDetails = false
				m.updateSelectionData(detailsDataUpdate)
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(details.Box)
				return m, cmd
			}

		case key.Matches(msg, keys.Mappings.Down):
			if m.stackTable.Focused() {
				m.stackTable.MoveDown(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.FocusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData(tasksDataUpdate)
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveDown(1)
				m.taskDetails.FocusIndex = 0
				m.showDetails = false
				m.updateSelectionData(detailsDataUpdate)
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(details.Box)
				return m, cmd
			}

		case key.Matches(msg, keys.Mappings.GotoTop):
			if m.stackTable.Focused() {
				m.stackTable.GotoTop()
				m.taskTable.SetCursor(0)
				m.taskDetails.FocusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData(tasksDataUpdate)
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoTop()
				m.taskDetails.FocusIndex = 0
				m.showDetails = false
				m.updateSelectionData(detailsDataUpdate)
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(details.Box)
				return m, cmd
			}

		case key.Matches(msg, keys.Mappings.GotoBottom):
			if m.stackTable.Focused() {
				m.stackTable.GotoBottom()
				m.taskTable.SetCursor(0)
				m.taskDetails.FocusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData(tasksDataUpdate)
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoBottom()
				m.taskDetails.FocusIndex = 0
				m.showDetails = false
				m.updateSelectionData(detailsDataUpdate)
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(details.Box)
				return m, cmd
			}

		case key.Matches(msg, keys.Mappings.New):
			if m.stackTable.Focused() {
				m.preInputFocus = stackViewName
				m.input = inputs.New(inputs.StackFormTable, entities.Stack{}, 0, m.context)

			} else if m.taskTable.Focused() {
				m.preInputFocus = taskViewName
				newTask := entities.Task{
					StackID: m.data[m.stackTable.Cursor()].ID,
				}
				m.input = inputs.New(inputs.TaskFormTable, newTask, 0, m.context)

			} else if m.taskDetails.Focused() {
				return m, nil
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.updateViewDimensions(14)

			m.showInput = true

			return m, nil

		case key.Matches(msg, keys.Mappings.Edit):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) == 0 {
					return m, nil
				}
				m.preInputFocus = stackViewName
				m.input = inputs.New(inputs.StackFormTable, m.data[m.stackTable.Cursor()], 0, m.context)
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = help.NewModel(keys.TaskDetailsMappings)
					m.navigationKeys = keys.DetailsMappings
				}
				return m, nil
			} else if m.taskDetails.Focused() {
				m.preInputFocus = detailViewName
				m.input = inputs.New(inputs.TaskFormTable, m.data[m.stackTable.Cursor()].Tasks[m.taskTable.Cursor()], m.taskDetails.FocusIndex, m.context)
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.updateViewDimensions(14)

			m.showInput = true

			return m, nil

		// Actual delete operation happens in showDelete conditional at the start of Update() method
		// Here we just trigger the delete confirmation step
		case key.Matches(msg, keys.Mappings.Delete):
			if m.stackTable.Focused() {
				m.preInputFocus = stackViewName
				m.showCustomInput = true
				m.customInputType = "delete"
				m.customInput = deletion.NewConfirmation()
				m.stackTable.Blur()
				m.help = help.Model{}

				return m, nil

			} else if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()

				if len(m.data[stackIndex].Tasks) > 0 {
					m.preInputFocus = taskViewName
					m.showCustomInput = true
					m.customInputType = "delete"
					m.customInput = deletion.NewConfirmation()
					m.taskTable.Blur()
					m.help = help.Model{}

					return m, nil
				}
			}

		case key.Matches(msg, keys.Mappings.Toggle):
			// Toggle task finish status
			if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()
				taskIndex := m.taskTable.Cursor()

				var currTask entities.Task
				if len(m.data[stackIndex].Tasks) > 0 {
					stack := m.data[stackIndex]
					currTask = stack.Tasks[taskIndex]

					// For recurring tasks we toggle the status of latest recur task entry
					currTask.IsFinished = !currTask.IsFinished
					currTask.Save(m.context)

					if currTask.IsFinished {
						stack.PendingTaskCount--
					} else {
						stack.PendingTaskCount++
					}
					stack.Save(m.context)

					stack.Tasks[taskIndex] = currTask
					m.data[stackIndex] = stack

					// Changing finish status will lead to reordering, so state has to be preserved
					m.preserveState()
					m.updateSelectionData(stacksDataUpdate)
					return m, nil
				}
			}

		case key.Matches(msg, keys.Mappings.Move):
			if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()

				if len(m.data[stackIndex].Tasks) > 0 {
					m.preInputFocus = taskViewName
					m.showCustomInput = true
					m.customInputType = "move"
					m.taskTable.Blur()

					opts := []definitions.KeyValue{}
					for _, stack := range m.data {
						entry := definitions.KeyValue{
							Key:   stack.ID,
							Value: stack.Title,
						}
						opts = append(opts, entry)
					}
					m.customInput = lists.NewSelector(opts, "", messages.MainGoToWith)

					m.help = help.NewModel(keys.ListSelectorMappings)
					return m, nil
				}
			}
		case key.Matches(msg, keys.Mappings.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, keys.Mappings.Quit, keys.Mappings.Exit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.context.Screen.Width = msg.Width
		m.context.Screen.Height = msg.Height
		m.updateViewDimensions(10)

		if m.firstRender {
			// updateSelectionData() is called here instead of being called from Init()
			// since details box rendering requires screen dimensions, which aren't set at the time of Init()
			m.updateSelectionData(stacksDataUpdate)
			m.firstRender = false
		}
	}

	return m, cmd
}

// View handles model view
func (m *model) View() string {
	var stackView, taskView, detailView string

	if m.stackTable.Focused() {
		stackView = display.SelectedStackBoxStyle.Render(m.stackView())
		taskView = display.UnselectedBoxStyle.Render(m.taskView())
		detailView = display.UnselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskTable.Focused() {
		stackView = display.UnselectedBoxStyle.Render(m.stackView())
		taskView = display.SelectedTaskBoxStyle.Render(m.taskView())
		detailView = display.UnselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskDetails.Focused() {
		stackView = display.UnselectedBoxStyle.Render(m.stackView())
		taskView = display.UnselectedBoxStyle.Render(m.taskView())
		detailView = display.SelectedDetailsBoxStyle.Render(m.taskDetails.View())
	} else {
		stackView = display.UnselectedBoxStyle.Render(m.stackView())
		taskView = display.UnselectedBoxStyle.Render(m.taskView())
		detailView = display.UnselectedBoxStyle.Render(m.taskDetails.View())
	}

	viewArr := []string{stackView}
	if m.showTasks {
		viewArr = append(viewArr, taskView)

		if m.showDetails {
			viewArr = append(viewArr, detailView)
		} else if len(m.taskTable.Rows()) > 0 {
			viewArr = append(viewArr, display.UnselectedBoxStyle.Render(m.context.Screen.EmptyDetailsView()))
		}
	} else {
		viewArr = append(viewArr, display.UnselectedBoxStyle.Render(m.context.Screen.EmptyTaskView()))
	}

	tablesView := lipgloss.JoinHorizontal(lipgloss.Center, viewArr...)

	if m.showCustomInput {
		tablesView = lipgloss.JoinVertical(lipgloss.Left,
			tablesView,
			m.context.Screen.InputFormStyle().Render(m.customInput.View()),
		)
	}

	if m.showInput {
		inputFormView := m.context.Screen.InputFormStyle().Render(m.input.View())
		tablesView = lipgloss.JoinVertical(lipgloss.Left,
			tablesView,
			inputFormView,
		)
		m.help = help.NewModel(m.input.HelpKeys())
	}

	if m.showHelp {
		if !m.showInput && !m.showCustomInput {
			navigationHelp := help.NewModel(m.navigationKeys)
			return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View(), navigationHelp.View())
		}
		return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View())
	}
	return tablesView
}

func (m *model) stackView() string {
	m.stackTable.SetHeight(m.context.Screen.Table.ViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.stackTable.View(), m.stackFooter())
}

func (m *model) stackFooter() string {
	stackFooterStyle := display.FooterContainerStyle.Width(display.StackTableWidth)

	info := display.FooterInfoStyle.Render(fmt.Sprintf("%d/%d", m.stackTable.Cursor()+1, len(m.stackTable.Rows())))

	return stackFooterStyle.Render(info)
}

func (m *model) taskView() string {
	m.taskTable.SetHeight(m.context.Screen.Table.ViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.taskTable.View(), m.taskFooter())
}

func (m *model) taskFooter() string {
	taskFooterStyle := display.FooterContainerStyle.Width(display.TaskTableWidth)

	if len(m.taskTable.Rows()) == 0 {
		return taskFooterStyle.Render("Press 'n' to create a new task")
	}
	info := display.FooterInfoStyle.Render(fmt.Sprintf("%d/%d", m.taskTable.Cursor()+1, len(m.taskTable.Rows())))
	return taskFooterStyle.Render(info)
}

// Pull new data from database
func (m *model) refreshData() {
	stacks, _ := entities.FetchStacks(m.context)
	m.data = stacks
	m.updateSelectionData(stacksDataUpdate)
}

// Efficiently update only the required pane
func (m *model) updateSelectionData(category string) {
	var retainIndex bool
	if m.prevState.retainState {
		retainIndex = true
		m.prevState.retainState = false
	}

	switch category {
	case stacksDataUpdate:
		m.updateStackTableData(retainIndex)
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(true)
	case tasksDataUpdate:
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(false)
	case detailsDataUpdate:
		m.updateDetailsBoxData(false)
	default:
		m.updateStackTableData(retainIndex)
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(true)
	}
}

func (m *model) updateStackTableData(retainIndex bool) {
	// Set stack view data
	// We pass a slice to stackRows, so the changes (like sorting) that happen there will be reflected in original slice
	m.stackTable.SetRows(tables.StackRows(m.data))

	if retainIndex {
		newIndex := entities.FindByIndex(m.data, m.prevState.stackID)

		if newIndex != -1 {
			m.stackTable.SetCursor(newIndex)
		}
	}
}

func (m *model) updateTaskTableData(retainIndex bool) {
	// Set task view data for selected stack
	stackIndex := m.stackTable.Cursor()
	currStack := m.data[stackIndex]

	// We pass a slice to taskRows, so the changes (like sorting) that happen there will be reflected in original slice
	m.taskTable.SetRows(tables.TaskRows(currStack.Tasks, m.context))

	if retainIndex {
		newIndex := entities.FindByIndex(m.data[stackIndex].Tasks, m.prevState.taskID)
		if newIndex != -1 {
			m.taskTable.SetCursor(newIndex)
		}
	}
}

func (m *model) updateDetailsBoxData(preserveOffset bool) {
	stackIndex := m.stackTable.Cursor()
	taskIndex := m.taskTable.Cursor()
	if taskIndex == -1 {
		taskIndex = 0
		m.taskTable.SetCursor(0)
	}

	var currTask entities.Task
	if len(m.data[stackIndex].Tasks) > 0 {
		currTask = m.data[stackIndex].Tasks[taskIndex]
	} else {
		currTask = entities.Task{}
	}

	m.taskDetails.Build(currTask, preserveOffset, m.context.Screen)
}

// Changing title, deadline, priority or finish status will lead to table reordering
// preserveState() is used to maintain focus on the stack/task that was being edited
func (m *model) preserveState() {
	m.prevState.retainState = true
	stackIndex := m.stackTable.Cursor()
	taskIndex := m.taskTable.Cursor()

	m.prevState.stackID = m.data[m.stackTable.Cursor()].ID
	if len(m.data[stackIndex].Tasks) > 0 {
		m.prevState.taskID = m.data[stackIndex].Tasks[taskIndex].ID
	}
}

func (m *model) updateViewDimensions(offset int) {
	m.context.Screen.Table.ViewHeight = m.context.Screen.Height - offset

	// Details box viewport dimensions & section width are set at the time of box creation,
	// after that they have to be manually adjusted
	m.taskDetails.ViewPort.Width = m.context.Screen.DetailsBoxWidth()
	m.taskDetails.ViewPort.Height = m.context.Screen.DetailsBoxHeight()
	m.updateDetailsBoxData(true)
}
