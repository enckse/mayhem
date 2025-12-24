package deletion_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/tui/deletion"
)

func TestConfirmation(t *testing.T) {
	obj := deletion.NewConfirmation()
	if obj.Init() == nil {
		t.Error("invalid result")
	}
	_, c := obj.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if c == nil {
		t.Error("invalid result")
	}
	_, c = obj.Update(tea.KeyMsg{Runes: []rune{'q'}})
	if c == nil {
		t.Error("invalid result")
	}
	_, c = obj.Update(tea.KeyMsg{Runes: []rune{'y'}})
	if c == nil {
		t.Error("invalid result")
	}
	v := obj.View()
	if !strings.Contains(v, "Do you wish to proceed with deletion") {
		t.Errorf("invalid deletion: %s", v)
	}
}
