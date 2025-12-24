package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
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
		taskDetails     detailsBox
		help            helpModel
		input           inputForm
		showTasks       bool
		showDetails     bool
		showInput       bool
		showHelp        bool
		customInput     tea.Model
		customInputType string
		showCustomInput bool
		navigationKeys  keyMap
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

// InitializeMainModel will startup the core application model
func InitializeMainModel(ctx *state.Context) ModelWrapper {
	stacks, _ := entities.FetchAllStacks(ctx)

	m := &model{
		stackTable:     buildTable(stackColumns(), "stack"),
		taskTable:      buildTable(taskColumns(), "task"),
		taskDetails:    detailsBox{}, // we can't build the details box at this stage since we need both stack & task indices for that
		data:           stacks,
		help:           initializeHelp(stackKeys),
		navigationKeys: tableNavigationKeys,
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

		case goToMainMsg:
			m.input = inputForm{}
			m.showInput = false

			if msg.value.(string) == "refresh" {
				m.preserveState()
				m.refreshData()
			}

			switch m.preInputFocus {
			case "stack":
				m.stackTable.Focus()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys
			case "task":
				m.taskTable.Focus()
				m.help = initializeHelp(taskKeys)
				m.navigationKeys = tableNavigationKeys
			case "detail":
				m.taskDetails.Focus()
				m.navigationKeys = detailsNavigationKeys
			}

			m.updateViewDimensions(10)

			return m, nil

		case tea.WindowSizeMsg:
			screenWidth = msg.Width
			screenHeight = msg.Height
			m.updateViewDimensions(14)
			return m, nil

		default:
			inp, cmd := m.input.Update(msg)
			t, _ := inp.(inputForm)
			m.input = t

			return m, cmd
		}
	}

	if m.showCustomInput {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			screenWidth = msg.Width
			screenHeight = msg.Height
			m.updateViewDimensions(14)
			return m, nil
		}

		switch m.customInputType {
		// Transfer control to delete confirmation model
		case "delete":
			switch msg := msg.(type) {

			case goToMainMsg:
				m.showCustomInput = false

				switch m.preInputFocus {
				case "stack":
					m.stackTable.Focus()
					m.help = initializeHelp(stackKeys)
				case "task":
					m.taskTable.Focus()
					m.help = initializeHelp(taskKeys)
				}

				if msg.value.(string) == "y" {
					switch m.preInputFocus {
					case "stack":
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

					case "task":
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
				t, _ := inp.(deleteConfirmation)
				m.customInput = t

				return m, cmd
			}

		case "move":
			switch msg := msg.(type) {

			case goToMainMsg:
				m.showCustomInput = false
				m.taskTable.Focus()
				m.help = initializeHelp(taskKeys)

				response := msg.value.(keyVal)

				if response.val == "" {
					return m, nil
				}

				newStackID := response.key

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
				entities.IncPendingCount(newStackID, m.context)
				currTask.StackID = newStackID
				currTask.Save(m.context)

				if taskIndex == len(m.taskTable.Rows())-1 {
					m.taskTable.SetCursor(taskIndex - 1)
				}
				m.refreshData()
				return m, nil

			default:
				inp, cmd := m.customInput.Update(msg)
				t, _ := inp.(listSelector)
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
		case key.Matches(msg, Keys.Left):
			if m.stackTable.Focused() {
				if m.showDetails {
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys

				}
			} else if m.taskTable.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys

			} else if m.taskDetails.Focused() {
				m.stackTable.Blur()
				m.taskTable.Focus()
				m.taskDetails.Blur()
				m.help = initializeHelp(taskKeys)
				m.navigationKeys = tableNavigationKeys

			}
			return m, nil

		case key.Matches(msg, Keys.Right):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) > 0 {
					m.showTasks = true
					m.stackTable.Blur()
					m.taskTable.Focus()
					m.taskDetails.Blur()
					m.help = initializeHelp(taskKeys)
					m.navigationKeys = tableNavigationKeys
					return m, nil
				}
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys
					return m, nil
				}
			} else if m.taskDetails.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys
				return m, nil
			}

		// Intra-table navigation

		// When we switch to a new stack:
		//		- Empty task box is shown
		//		- Details box is hidden

		// When we switch to a new task:
		//		- Empty details box is shown
		case key.Matches(msg, Keys.Up):
			if m.stackTable.Focused() {
				m.stackTable.MoveUp(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveUp(1)
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.Down):
			if m.stackTable.Focused() {
				m.stackTable.MoveDown(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveDown(1)
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.GotoTop):
			if m.stackTable.Focused() {
				m.stackTable.GotoTop()
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoTop()
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.GotoBottom):
			if m.stackTable.Focused() {
				m.stackTable.GotoBottom()
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoBottom()
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.New):
			if m.stackTable.Focused() {
				m.preInputFocus = "stack"
				m.input = initializeInput("stack", entities.Stack{}, 0, m.context)

			} else if m.taskTable.Focused() {
				m.preInputFocus = "task"
				newTask := entities.Task{
					StackID: m.data[m.stackTable.Cursor()].ID,
				}
				m.input = initializeInput("task", newTask, 0, m.context)

			} else if m.taskDetails.Focused() {
				return m, nil
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.updateViewDimensions(14)

			m.showInput = true

			return m, nil

		case key.Matches(msg, Keys.Edit):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) == 0 {
					return m, nil
				}
				m.preInputFocus = "stack"
				m.input = initializeInput("stack", m.data[m.stackTable.Cursor()], 0, m.context)
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys
				}
				return m, nil
			} else if m.taskDetails.Focused() {
				m.preInputFocus = "detail"
				m.input = initializeInput("task", m.data[m.stackTable.Cursor()].Tasks[m.taskTable.Cursor()], m.taskDetails.focusIndex, m.context)
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.updateViewDimensions(14)

			m.showInput = true

			return m, nil

		// Actual delete operation happens in showDelete conditional at the start of Update() method
		// Here we just trigger the delete confirmation step
		case key.Matches(msg, Keys.Delete):
			if m.stackTable.Focused() {
				m.preInputFocus = "stack"
				m.showCustomInput = true
				m.customInputType = "delete"
				m.customInput = initializeDeleteConfirmation()
				m.stackTable.Blur()
				m.help = helpModel{}

				return m, nil

			} else if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()

				if len(m.data[stackIndex].Tasks) > 0 {
					m.preInputFocus = "task"
					m.showCustomInput = true
					m.customInputType = "delete"
					m.customInput = initializeDeleteConfirmation()
					m.taskTable.Blur()
					m.help = helpModel{}

					return m, nil
				}
			}

		case key.Matches(msg, Keys.Toggle):
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
					m.updateSelectionData("stacks")
					return m, nil
				}
			}

		case key.Matches(msg, Keys.Move):
			if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()

				if len(m.data[stackIndex].Tasks) > 0 {
					m.preInputFocus = "task"
					m.showCustomInput = true
					m.customInputType = "move"
					m.taskTable.Blur()

					opts := []keyVal{}
					for _, stack := range m.data {
						entry := keyVal{
							key: stack.ID,
							val: stack.Title,
						}
						opts = append(opts, entry)
					}
					m.customInput = initializeListSelector(opts, "", goToMainWithVal)

					m.help = initializeHelp(listSelectorKeys)
					return m, nil
				}
			}
		case key.Matches(msg, Keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, Keys.Quit, Keys.Exit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		screenWidth = msg.Width
		screenHeight = msg.Height
		m.updateViewDimensions(10)

		if m.firstRender {
			// updateSelectionData() is called here instead of being called from Init()
			// since details box rendering requires screen dimensions, which aren't set at the time of Init()
			m.updateSelectionData("stacks")
			m.firstRender = false
		}
	}

	return m, cmd
}

// View handles model view
func (m *model) View() string {
	var stackView, taskView, detailView string

	if m.stackTable.Focused() {
		stackView = selectedStackBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskTable.Focused() {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = selectedTaskBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskDetails.Focused() {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = selectedDetailsBoxStyle.Render(m.taskDetails.View())
	} else {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	}

	viewArr := []string{stackView}
	if m.showTasks {
		viewArr = append(viewArr, taskView)

		if m.showDetails {
			viewArr = append(viewArr, detailView)
		} else if len(m.taskTable.Rows()) > 0 {
			viewArr = append(viewArr, unselectedBoxStyle.Render(getEmptyDetailsView()))
		}
	} else {
		viewArr = append(viewArr, unselectedBoxStyle.Render(getEmptyTaskView()))
	}

	tablesView := lipgloss.JoinHorizontal(lipgloss.Center, viewArr...)

	if m.showCustomInput {
		tablesView = lipgloss.JoinVertical(lipgloss.Left,
			tablesView,
			getInputFormStyle().Render(m.customInput.View()),
		)
	}

	if m.showInput {
		inputFormView := getInputFormStyle().Render(m.input.View())
		tablesView = lipgloss.JoinVertical(lipgloss.Left,
			tablesView,
			inputFormView,
		)
		m.help = initializeHelp(m.input.helpKeys)
	}

	if m.showHelp {
		if !m.showInput && !m.showCustomInput {
			navigationHelp := initializeHelp(m.navigationKeys)
			return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View(), navigationHelp.View())
		}
		return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View())
	}
	return tablesView
}

func (m *model) stackView() string {
	m.stackTable.SetHeight(tableViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.stackTable.View(), m.stackFooter())
}

func (m *model) stackFooter() string {
	stackFooterStyle := footerContainerStyle.Width(stackTableWidth)

	info := footerInfoStyle.Render(fmt.Sprintf("%d/%d", m.stackTable.Cursor()+1, len(m.stackTable.Rows())))

	return stackFooterStyle.Render(info)
}

func (m *model) taskView() string {
	m.taskTable.SetHeight(tableViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.taskTable.View(), m.taskFooter())
}

func (m *model) taskFooter() string {
	taskFooterStyle := footerContainerStyle.Width(taskTableWidth)

	if len(m.taskTable.Rows()) == 0 {
		return taskFooterStyle.Render("Press 'n' to create a new task")
	}
	info := footerInfoStyle.Render(fmt.Sprintf("%d/%d", m.taskTable.Cursor()+1, len(m.taskTable.Rows())))
	return taskFooterStyle.Render(info)
}

// Pull new data from database
func (m *model) refreshData() {
	stacks, _ := entities.FetchAllStacks(m.context)
	m.data = stacks
	m.updateSelectionData("stacks")
}

// Efficiently update only the required pane
func (m *model) updateSelectionData(category string) {
	var retainIndex bool
	if m.prevState.retainState {
		retainIndex = true
		m.prevState.retainState = false
	}

	switch category {
	case "stacks":
		m.updateStackTableData(retainIndex)
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(true)
	case "tasks":
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(false)
	case "details":
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
	m.stackTable.SetRows(stackRows(m.data))

	if retainIndex {
		newIndex := findIndex(m.data, m.prevState.stackID)

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
	m.taskTable.SetRows(taskRows(currStack.Tasks))

	if retainIndex {
		newIndex := findIndex(m.data[stackIndex].Tasks, m.prevState.taskID)
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

	m.taskDetails.buildDetailsBox(currTask, preserveOffset)
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
	tableViewHeight = screenHeight - offset

	// Details box viewport dimensions & section width are set at the time of box creation,
	// after that they have to be manually adjusted
	m.taskDetails.viewport.Width = getDetailsBoxWidth()
	m.taskDetails.viewport.Height = getDetailsBoxHeight()
	m.updateDetailsBoxData(true)
}
