package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/tui/display"
	"github.com/enckse/mayhem/internal/tui/keys"
)

type (
	detailsBox struct {
		taskData          entities.Task
		viewport          viewport.Model
		preserveOffset    bool
		oldViewportOffset int
		focusIndex        int
		isBoxFocused      bool
		scrollData        scrollData
	}

	scrollData struct {
		title    int
		notes    int
		priority int
		deadline int
	}
)

var (
	taskDetailsKeys = keys.Map{
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit field 📝"),
		),
	}

	detailsNavigationKeys = keys.Map{
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
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("'?'", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("'q'", "quit"),
		),
	}
)

func (m *detailsBox) buildDetailsBox(data entities.Task, preserveOffset bool) {
	m.taskData = data

	// We want to preserve offset when we return to same details view after editing any field
	// But when going from one task to another, we want to reset the view
	m.preserveOffset = preserveOffset
	m.oldViewportOffset = m.viewport.YOffset
	m.viewport = viewport.New(display.DetailsBoxWidth(), display.TableViewHeight)
	m.renderContent()
}

func (m detailsBox) Init() tea.Cmd {
	return nil
}

func (m detailsBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isBoxFocused {
		return m, nil
	}

	m.viewport.Width = display.DetailsBoxWidth()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Mappings.Up):
			var scrollDistance int
			switch m.focusIndex {
			case taskTitleIndex:
				m.viewport.GotoBottom()
				m.End()
				return m, nil
			case taskNotesIndex:
				scrollDistance = m.scrollData.notes
				m.Previous()
			case taskPriorityIndex:
				scrollDistance = m.scrollData.priority
				m.Previous()
			case taskDeadlineIndex:
				m.Previous()
			}

			m.viewport.ScrollUp(scrollDistance)

		case key.Matches(msg, keys.Mappings.Down):
			var scrollDistance int
			switch m.focusIndex {
			case taskTitleIndex:
				m.Next()
			case taskNotesIndex:
				scrollDistance = m.scrollData.notes
				m.Next()
			case taskPriorityIndex:
				scrollDistance = m.scrollData.priority
				m.Next()
			case taskDeadlineIndex:
				m.viewport.GotoTop()
				m.Start()
				return m, nil
			}

			m.viewport.ScrollDown(scrollDistance)

		case key.Matches(msg, keys.Mappings.GotoTop):
			m.viewport.GotoTop()
			m.Start()

		case key.Matches(msg, keys.Mappings.GotoBottom):
			m.viewport.GotoBottom()
			m.End()
		}
	}
	return m, nil
}

func (m detailsBox) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, display.DetailsBoxStyle().Render(m.viewport.View()), m.footerView())
}

func (m *detailsBox) Focus() {
	m.isBoxFocused = true
}

func (m *detailsBox) Blur() {
	m.isBoxFocused = false
}

func (m detailsBox) Focused() bool {
	return m.isBoxFocused
}

func (m *detailsBox) Next() {
	length := taskLastIndex + 1
	m.focusIndex = (m.focusIndex + 1) % length
	m.renderContent()
}

func (m *detailsBox) End() {
	m.focusIndex = taskLastIndex
	m.renderContent()
}

func (m *detailsBox) Previous() {
	length := taskLastIndex + 1
	val := (m.focusIndex - 1) % length
	if val < 0 {
		val = val + length
	}
	m.focusIndex = val
	m.renderContent()
}

func (m *detailsBox) Start() {
	m.focusIndex = 0
	m.renderContent()
}

func (m *detailsBox) renderContent() {
	content := []string{
		m.titleBlock(),
		m.notesBlock(),
		m.priorityBlock(),
		m.deadlineBlock(),
	}

	view := lipgloss.JoinVertical(lipgloss.Left, content...)
	m.viewport.SetContent(view)
	if m.preserveOffset {
		m.viewport.SetYOffset(m.oldViewportOffset)
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

func (m *detailsBox) titleBlock() string {
	var b strings.Builder
	isFocused := (m.focusIndex == taskTitleIndex)
	newBlock(&b, "Title", isFocused)
	b.WriteString(m.taskData.Title)

	data := display.ItemContainerStyle(isFocused).Render(display.DetailsItemStyle(isFocused).PaddingTop(0).Render(b.String()))
	m.scrollData.title = lipgloss.Height(data)
	return data
}

func (m *detailsBox) notesBlock() string {
	var b strings.Builder
	isFocused := (m.focusIndex == taskNotesIndex)
	newBlock(&b, "Notes", isFocused)

	if m.taskData.Notes == "" {
		b.WriteString("-")
	} else {
		b.WriteString(m.taskData.Notes)
	}

	data := display.ItemContainerStyle(isFocused).Render(display.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.notes = lipgloss.Height(data)
	return data
}

func (m *detailsBox) priorityBlock() string {
	var b strings.Builder
	isFocused := (m.focusIndex == taskPriorityIndex)
	newBlock(&b, "Priority", isFocused)
	fmt.Fprintf(&b, "%d", m.taskData.Priority)

	data := display.ItemContainerStyle(isFocused).Render(display.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.priority = lipgloss.Height(data)
	return data
}

func (m *detailsBox) deadlineBlock() string {
	var b strings.Builder
	isFocused := (m.focusIndex == taskDeadlineIndex)
	newBlock(&b, "Deadline", isFocused)

	if m.taskData.Deadline.IsZero() {
		b.WriteString("Not Scheduled")
	} else {
		b.WriteString(formatTime(m.taskData.Deadline, true))
	}

	data := display.ItemContainerStyle(isFocused).Render(display.DetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.deadline = lipgloss.Height(data)
	return data
}

func (m *detailsBox) footerView() string {
	scrollInfoStyle := display.FooterContainerStyle.Width(m.viewport.Width).Align(lipgloss.Right)
	info := display.FooterInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	return scrollInfoStyle.Render(info)
}
