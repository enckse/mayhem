// Package display deals with styles and screen management definitions
package display

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// TableType are the various UI display tables that exist
type TableType uint

const (
	// StackTableWidth is the width of the stack(s) table
	// 22: column width + 2*2: column padding
	StackTableWidth = 26
	// TaskTableWidth is the width for the actual task list table
	// 59: column widths + 2*4: column paddings
	TaskTableWidth = 67
	// StackTableType defines the stack table definition
	StackTableType TableType = iota
	// TaskTableType defines the task table definition
	TaskTableType
)

var (
	// TableViewHeight is the height of displayable tables
	TableViewHeight = 25
	// ScreenWidth is the width of the screen
	ScreenWidth int
	// ScreenHeight is the height of the screen
	ScreenHeight int

	inputFormBorderColor = lipgloss.Color("#325b84")

	taskSelectionColor    = lipgloss.Color("#f1b44c")
	detailsSelectionColor = lipgloss.Color("#333c4d")

	// HighlightedBackgroundColor indicates a highlighted background
	HighlightedBackgroundColor = lipgloss.Color("#f97171")
	highlightedTextColor       = lipgloss.Color("#4e4e4e")
	// InputFormColor is the color for input form(s)
	InputFormColor = lipgloss.Color("#5ac7c7")
	// TimeFocusColor is used to show time control focus
	TimeFocusColor = lipgloss.Color("#FFFF00")
	// UnfocusedColor is a general unfocus indicator
	UnfocusedColor = lipgloss.Color("#898989")

	selectedBoxStyle = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder())
	// SelectedStackBoxStyle indicates the stack box is focused
	SelectedStackBoxStyle = selectedBoxStyle.BorderForeground(lipgloss.Color("#019187"))
	// SelectedTaskBoxStyle indicates the task box is focused
	SelectedTaskBoxStyle = selectedBoxStyle.BorderForeground(lipgloss.Color("#f1b44c"))
	// SelectedDetailsBoxStyle indicates the details box is focused
	SelectedDetailsBoxStyle = selectedBoxStyle.BorderForeground(lipgloss.Color("#6192bc"))
	// UnselectedBoxStyle are for the other (unfocused) boxes
	UnselectedBoxStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(UnfocusedColor)
	stackSelectedRowStyle = table.DefaultStyles().Selected.
				Foreground(highlightedTextColor).
				Background(lipgloss.Color("#019187")).
				Bold(false)
	taskSelectedRowStyle = stackSelectedRowStyle.Background(taskSelectionColor)
	// FooterInfoStyle indicates how the footer is styled
	FooterInfoStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#1c2c4c"))
	// FooterContainerStyle is the overall container for the footer
	FooterContainerStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Background(lipgloss.Color("#3e424b"))
	// HighlightedTextStyle is the style for highlighting text
	HighlightedTextStyle = lipgloss.NewStyle().
				Bold(true).
				Italic(true).
				Foreground(highlightedTextColor).
				Background(HighlightedBackgroundColor).
				Padding(0, 1).
				MarginTop(1)
	// TextInputStyle handles text inputs
	TextInputStyle = lipgloss.NewStyle().Foreground(InputFormColor)
	// PlaceHolderStyle is for placeholder styles
	PlaceHolderStyle = lipgloss.NewStyle().Foreground(UnfocusedColor)
)

// InputFormStyle will get the default input form style to use
// Since width is dynamic, we have to append it to the style before usage
func InputFormStyle() lipgloss.Style {
	// Subtract 2 for padding on each side
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(inputFormBorderColor).
		Padding(0, 1).Width(ScreenWidth - 2)
}

// TableStyle will get the table style definition
func TableStyle(tableType TableType) table.Styles {
	s := table.DefaultStyles()
	s.Header = table.DefaultStyles().Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(UnfocusedColor).
		BorderBottom(true).
		Bold(true)

	switch tableType {
	case StackTableType:
		s.Selected = stackSelectedRowStyle
	case TaskTableType:
		s.Selected = taskSelectedRowStyle
	}

	return s
}

// EmptyTaskStyle will get the style for an empty task
func EmptyTaskStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(60).
		Height(TableViewHeight + 3) // 3 is added to account for header & footer height
}

// EmptyDetailsStyle will get the style for empty details
func EmptyDetailsStyle() lipgloss.Style {
	return DetailsBoxStyle().
		Height(TableViewHeight + 3).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
}

// DetailsBoxWidth will get the width for the details box
func DetailsBoxWidth() int {
	return ScreenWidth - (StackTableWidth + TaskTableWidth + 3*2) // each of the 3 boxes have left & right borders
}

// DetailsBoxHeight will get the height for the details box
func DetailsBoxHeight() int {
	return TableViewHeight
}

// DetailsBoxStyle will get the style for details box
func DetailsBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(DetailsBoxWidth()).
		Height(TableViewHeight + 2)
}

// DetailsItemStyle will get the style for details box item(s)
func DetailsItemStyle(isSelected bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(0, 0, 1, 0).
		Width(DetailsBoxWidth() - 2)

	if isSelected {
		style.Background(detailsSelectionColor)
	}

	return style
}

// ItemContainerStyle will get the style for an item based on default style and selection
// Applying padding (0,1) to detail items causes issue with description text alignment
// To avoid that an additional container is used for detail items
func ItemContainerStyle(isSelected bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(DetailsBoxWidth())

	if isSelected {
		style.Background(detailsSelectionColor)
	}

	return style
}
