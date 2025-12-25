package textarea_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/inputs/textarea"
	"github.com/enckse/mayhem/internal/tui/messages"
)

func TestInput(t *testing.T) {
	obj := textarea.New("", &display.Screen{})
	if obj.Init() == nil {
		t.Error("invalid result")
	}
	_, c := obj.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if c == nil {
		t.Error("invalid result")
	}
	if _, ok := c().(messages.Form); !ok {
		t.Error("invalid command")
	}
	_, c = obj.Update(tea.KeyMsg{Runes: []rune{'c'}})
	if c == nil {
		t.Error("invalid result")
	}
	if _, ok := c().(messages.Form); ok {
		t.Error("invalid command")
	}
	v := obj.View()
	if !strings.Contains(v, "┃") {
		t.Errorf("invalid input: %s", v)
	}
}
