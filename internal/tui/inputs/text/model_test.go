package text_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/tui/inputs/text"
	"github.com/enckse/mayhem/internal/tui/messages"
)

func TestInput(t *testing.T) {
	obj := text.New("", "", 100, messages.FormGoToWith)
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
	if !strings.Contains(v, ">") {
		t.Errorf("invalid input: %s", v)
	}
}
