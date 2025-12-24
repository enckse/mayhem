package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/tui/keys"
)

type (
	// textinput.Model doesn't implement tea.Model interface
	timePicker struct {
		currTime   time.Time
		focusIndex int
	}

	timeUnit struct {
		title     string
		tag       string
		charWidth int
	}
)

const (
	dayItem int = iota
	monthItem
	yearItem
	hourItem
	minuteItem
)

var (
	timePickerKeys = keys.Map{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("'↑/k'", "increase"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("'↓/j'", "decrease"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("'←/h'", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("'→/l'", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("'enter'", "save"),
		),
		Return: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("'esc'", "return"),
		),
	}

	timeUnitMap = map[int]timeUnit{
		hourItem: {
			title:     "Hour",
			tag:       "hh",
			charWidth: 2,
		},
		minuteItem: {
			title:     "Minute",
			tag:       "mm",
			charWidth: 2,
		},
		dayItem: {
			title:     "Day",
			tag:       "DD",
			charWidth: 2,
		},
		monthItem: {
			title:     "Month",
			tag:       "MM",
			charWidth: 2,
		},
		yearItem: {
			title:     "Year",
			tag:       "YYYY",
			charWidth: 4,
		},
	}
)

func initializeTimePicker(currTime time.Time) tea.Model {
	t := timePicker{
		currTime: currTime,
	}

	return t
}

func (m timePicker) Init() tea.Cmd {
	return nil
}

func (m timePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Mappings.Up):
			switch m.focusIndex {
			case hourItem:
				hourDuration, _ := time.ParseDuration("60m")
				m.currTime = m.currTime.Add(hourDuration)
			case minuteItem:
				minuteDuration, _ := time.ParseDuration("1m")
				m.currTime = m.currTime.Add(minuteDuration)
			case dayItem:
				m.currTime = m.currTime.AddDate(0, 0, 1)
			case monthItem:
				m.currTime = m.currTime.AddDate(0, 1, 0)
			case yearItem:
				m.currTime = m.currTime.AddDate(1, 0, 0)
			}
			return m, nil

		case key.Matches(msg, keys.Mappings.Down):
			switch m.focusIndex {
			case hourItem:
				hourDuration, _ := time.ParseDuration("60m")
				m.currTime = m.currTime.Add(-hourDuration)
			case minuteItem:
				minuteDuration, _ := time.ParseDuration("1m")
				m.currTime = m.currTime.Add(-minuteDuration)
			case dayItem:
				m.currTime = m.currTime.AddDate(0, 0, -1)
			case monthItem:
				m.currTime = m.currTime.AddDate(0, -1, 0)
			case yearItem:
				m.currTime = m.currTime.AddDate(-1, 0, 0)
			}
			return m, nil

		case key.Matches(msg, keys.Mappings.Right):
			if m.focusIndex < len(timeUnitMap)-1 {
				m.focusIndex++
			}
			return m, nil

		case key.Matches(msg, keys.Mappings.Left):
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			return m, nil
		case key.Matches(msg, keys.Mappings.Enter):
			return m, goToFormWithVal(m.currTime)
		}
	}
	return m, nil
}

func (m timePicker) View() string {
	var timeUnitLabel string
	var timeValue string

	// Empty spaces are added to align the label and value rows
	timeUnitLabel = lipgloss.JoinHorizontal(lipgloss.Center,
		m.renderUnitTag(dayItem),
		" ",
		m.renderUnitTag(monthItem),
		" ",
		m.renderUnitTag(yearItem),
		" ",
		m.renderUnitTag(hourItem),
		" ",
		m.renderUnitTag(minuteItem),
	)

	timeValue = lipgloss.JoinHorizontal(lipgloss.Center,
		m.renderUnitCol(dayItem, m.currTime.Day()),
		"-",
		m.renderUnitCol(monthItem, int(m.currTime.Month())),
		"-",
		m.renderUnitCol(yearItem, m.currTime.Year()),
		" ",
		m.renderUnitCol(hourItem, formatHour(m.currTime.Hour())),
		":",
		m.renderUnitCol(minuteItem, m.currTime.Minute()),
		" ",
		renderMidDayInfo(m.currTime.Hour()))

	return lipgloss.JoinVertical(lipgloss.Center,
		timeValue,
		timeUnitLabel,
	)
}

func (m timePicker) renderUnitCol(index, val int) string {
	value := fmt.Sprintf("%0*d", timeUnitMap[index].charWidth, val)

	var color lipgloss.Color
	if m.focusIndex == index {
		color = timeFocusColor
	} else {
		color = unfocusedColor
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		BorderForeground(color).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	return style.Render(value)
}

func (m timePicker) renderUnitTag(index int) string {
	value := timeUnitMap[index].tag

	var color lipgloss.Color
	if m.focusIndex == index {
		color = timeFocusColor
	} else {
		color = unfocusedColor
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Padding(0, 2)

	return style.Render(value)
}

func renderMidDayInfo(hours int) string {
	if hours >= 12 {
		return "pm"
	}
	return "am"
}

// Adjust Hour value to 12 hour clock format
func formatHour(value int) int {
	if value > 12 {
		return value - 12
	}
	return value
}
