// Package details will show detail UI element
package details

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/inputs/timepicker"
	"github.com/enckse/mayhem/internal/tui/keys"
)

type (
	// Box is the details box
	Box struct {
		taskData          entities.Task
		ViewPort          viewport.Model
		preserveOffset    bool
		oldViewportOffset int
		FocusIndex        int
		isBoxFocused      bool
		scrollData        scrollData
		screen            *display.Screen
	}

	scrollData struct {
		title    int
		notes    int
		priority int
		deadline int
	}
)

// NewBox will create a new details box
func NewBox(screen *display.Screen) Box {
	return Box{screen: screen}
}

// Build will construct a new details box
func (m *Box) Build(data entities.Task, preserveOffset bool) {
	m.taskData = data

	// We want to preserve offset when we return to same details view after editing any field
	// But when going from one task to another, we want to reset the view
	m.preserveOffset = preserveOffset
	m.oldViewportOffset = m.ViewPort.YOffset
	m.ViewPort = viewport.New(m.screen.DetailsBoxWidth(), m.screen.Table.ViewHeight)
	m.renderContent()
}

// Init will init the model
func (m Box) Init() tea.Cmd {
	return nil
}

// Update will update the model
func (m Box) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isBoxFocused {
		return m, nil
	}

	m.ViewPort.Width = m.screen.DetailsBoxWidth()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Mappings.Up):
			var scrollDistance int
			switch m.FocusIndex {
			case definitions.TaskTitleIndex:
				m.ViewPort.GotoBottom()
				m.End()
				return m, nil
			case definitions.TaskNotesIndex:
				scrollDistance = m.scrollData.notes
				m.Previous()
			case definitions.TaskPriorityIndex:
				scrollDistance = m.scrollData.priority
				m.Previous()
			case definitions.TaskDeadlineIndex:
				m.Previous()
			}

			m.ViewPort.ScrollUp(scrollDistance)

		case key.Matches(msg, keys.Mappings.Down):
			var scrollDistance int
			switch m.FocusIndex {
			case definitions.TaskTitleIndex:
				m.Next()
			case definitions.TaskNotesIndex:
				scrollDistance = m.scrollData.notes
				m.Next()
			case definitions.TaskPriorityIndex:
				scrollDistance = m.scrollData.priority
				m.Next()
			case definitions.TaskDeadlineIndex:
				m.ViewPort.GotoTop()
				m.Start()
				return m, nil
			}

			m.ViewPort.ScrollDown(scrollDistance)

		case key.Matches(msg, keys.Mappings.GotoTop):
			m.ViewPort.GotoTop()
			m.Start()

		case key.Matches(msg, keys.Mappings.GotoBottom):
			m.ViewPort.GotoBottom()
			m.End()
		}
	}
	return m, nil
}

// View will display the model
func (m Box) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, m.screen.DetailsBoxStyle().Render(m.ViewPort.View()), m.footerView())
}

// Focus will set box focus
func (m *Box) Focus() {
	m.isBoxFocused = true
}

// Blur will set box blur
func (m *Box) Blur() {
	m.isBoxFocused = false
}

// Focused will indicate if focused
func (m Box) Focused() bool {
	return m.isBoxFocused
}

// Next will move to the next component
func (m *Box) Next() {
	length := definitions.TaskLastIndex + 1
	m.FocusIndex = (m.FocusIndex + 1) % length
	m.renderContent()
}

// End will move the end
func (m *Box) End() {
	m.FocusIndex = definitions.TaskLastIndex
	m.renderContent()
}

// Previous will move to the previous component
func (m *Box) Previous() {
	length := definitions.TaskLastIndex + 1
	val := (m.FocusIndex - 1) % length
	if val < 0 {
		val = val + length
	}
	m.FocusIndex = val
	m.renderContent()
}

// Start will move to the first component
func (m *Box) Start() {
	m.FocusIndex = 0
	m.renderContent()
}

func (m *Box) renderContent() {
	content := []string{
		m.titleBlock(),
		m.notesBlock(),
		m.priorityBlock(),
		m.deadlineBlock(),
	}

	view := lipgloss.JoinVertical(lipgloss.Left, content...)
	m.ViewPort.SetContent(view)
	if m.preserveOffset {
		m.ViewPort.SetYOffset(m.oldViewportOffset)
		m.preserveOffset = false
	}
}

func newBlock(b *strings.Builder, title string, isFocus bool) {
	prefix := ""
	if isFocus {
		prefix = "» "
	}
	b.WriteString(display.HighlightedTextStyle.Render(fmt.Sprintf("%s%s:", prefix, title)))
	b.WriteString("\n\n")
}

func (m *Box) titleBlock() string {
	var b strings.Builder
	isFocused := (m.FocusIndex == definitions.TaskTitleIndex)
	newBlock(&b, "Title", isFocused)
	b.WriteString(m.taskData.Title)

	data := m.screen.ItemContainerStyle(isFocused).Render(m.screen.DetailsItemStyle(isFocused).PaddingTop(0).Render(b.String()))
	m.scrollData.title = lipgloss.Height(data)
	return data
}

func (m *Box) notesBlock() string {
	var b strings.Builder
	isFocused := (m.FocusIndex == definitions.TaskNotesIndex)
	newBlock(&b, "Notes", isFocused)

	if m.taskData.Notes == "" {
		b.WriteString("-")
	} else {
		b.WriteString(m.taskData.Notes)
	}

	data := m.screen.ItemContainerStyle(isFocused).Render(m.screen.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.notes = lipgloss.Height(data)
	return data
}

func (m *Box) priorityBlock() string {
	var b strings.Builder
	isFocused := (m.FocusIndex == definitions.TaskPriorityIndex)
	newBlock(&b, "Priority", isFocused)
	fmt.Fprintf(&b, "%d", m.taskData.Priority)

	data := m.screen.ItemContainerStyle(isFocused).Render(m.screen.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.priority = lipgloss.Height(data)
	return data
}

func (m *Box) deadlineBlock() string {
	var b strings.Builder
	isFocused := (m.FocusIndex == definitions.TaskDeadlineIndex)
	newBlock(&b, "Deadline", isFocused)

	if m.taskData.Deadline.IsZero() {
		b.WriteString("Not Scheduled")
	} else {
		b.WriteString(timepicker.FormatTime(m.taskData.Deadline, true))
	}

	data := m.screen.ItemContainerStyle(isFocused).Render(m.screen.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.deadline = lipgloss.Height(data)
	return data
}

func (m *Box) footerView() string {
	scrollInfoStyle := display.FooterContainerStyle.Width(m.ViewPort.Width).Align(lipgloss.Right)
	info := display.FooterInfoStyle.Render(fmt.Sprintf("%3.f%%", m.ViewPort.ScrollPercent()*100))
	return scrollInfoStyle.Render(info)
}
