package lists_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/inputs/lists"
	"github.com/enckse/mayhem/internal/tui/messages"
)

func TestInput(t *testing.T) {
	obj := lists.NewSelector([]definitions.KeyValue{{Value: "abc"}, {Value: "xyz"}}, "abc", messages.MainGoToWith)
	if !strings.Contains(obj.View(), "» abc") {
		t.Errorf("invalid view: %s", obj.View())
	}
	_, cmd := obj.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if _, ok := cmd().(messages.Main); !ok {
		t.Error("invalid result")
	}
	_, cmd = obj.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("invalid result")
	}
	_, cmd = obj.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if cmd == nil {
		t.Error("invalid result")
	}
	_, cmd = obj.Update(tea.KeyMsg{Type: tea.KeyDown})
	if cmd != nil {
		t.Error("invalid result")
	}
	_, cmd = obj.Update(tea.KeyMsg{Type: tea.KeyUp})
	if cmd != nil {
		t.Error("invalid result")
	}
	_, cmd = obj.Update("")
	if cmd != nil {
		t.Error("invalid result")
	}
}
