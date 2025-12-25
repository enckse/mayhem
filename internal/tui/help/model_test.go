package help_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/tui/help"
	"github.com/enckse/mayhem/internal/tui/keys"
)

func TestModel(t *testing.T) {
	m := help.NewModel(keys.Mappings)
	if m.Init() != nil {
		t.Error("invalid init")
	}
	m.Update(tea.WindowSizeMsg{Width: 1})
	if !strings.Contains(m.View(), "'enter' enter") {
		t.Error("view failed")
	}
}
