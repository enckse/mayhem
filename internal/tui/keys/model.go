// Package keys handles key mappings
package keys

import "github.com/charmbracelet/bubbles/key"

// Map is the key binding map definition
type Map struct {
	Up            key.Binding
	Down          key.Binding
	GotoTop       key.Binding
	GotoBottom    key.Binding
	Left          key.Binding
	Right         key.Binding
	New           key.Binding
	Edit          key.Binding
	Move          key.Binding
	Enter         key.Binding
	Save          key.Binding
	Toggle        key.Binding
	ReverseToggle key.Binding
	Delete        key.Binding
	Return        key.Binding
	Help          key.Binding
	Quit          key.Binding
	Exit          key.Binding
}

var (
	// Mappings are actual key bindings across the app
	Mappings = Map{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("'↑/k'", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("'↓/j'", "move down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("'g'", "go to top"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("'G'", "go to bottom"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("'←/h'", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("'→/l'", "move right"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("'n'", "new"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("'m'", "move"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("'enter'", "enter"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("'tab'", "toggle"),
		),
		ReverseToggle: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("'shift+tab'", "toggle"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("'x'", "delete 🗑"),
		),
		Return: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("'esc'", "return"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("'?'", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("'q'", "quit"),
		),
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("'ctrl+c'", "exit"),
		),
	}

	// TextInputMappings are for form text fields
	TextInputMappings = Map{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("'enter'", "save"),
		),
		Return: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("'esc'", "return"),
		),
	}
	// ListSelectorMappings are for list selections
	ListSelectorMappings = Map{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("'↑/k'", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("'↓/j'", "down"),
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
	// TimePickerMappings are for the time picker
	TimePickerMappings = Map{
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
	// TextAreaInputMappings are for text areas
	TextAreaInputMappings = Map{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("'enter'", "new line"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("'ctrl+s'", "save"),
		),
		Return: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("'esc'", "return"),
		),
	}
	// DetailsMappings handle moving through the details screen
	DetailsMappings = Map{
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
	// TaskDetailsMappings manage editing a task
	TaskDetailsMappings = Map{
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit field"),
		),
	}

	// StackMappings navigate the stack
	StackMappings = Map{
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("'n'", "new stack"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("'x'", "delete"),
		),
	}

	// TaskMappings navigate the tasks
	TaskMappings = Map{
		Toggle: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("'tab'", "check/uncheck"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("'n'", "new task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("'e'", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("'x'", "delete"),
		),
		Move: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("'m'", "change stack"),
		),
	}

	// TableMappings navigate a table
	TableMappings = Map{
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
)

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k Map) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Toggle,
		k.ReverseToggle,
		k.New,
		k.Edit,
		k.Enter,
		k.Save,
		k.Delete,
		k.Move,
		k.Return,
		k.Up,
		k.Down,
		k.GotoTop,
		k.GotoBottom,
		k.Left,
		k.Right,
		k.Help,
		k.Quit,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k Map) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
