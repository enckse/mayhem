package details_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/tui/details"
)

func TestBuild(t *testing.T) {
	b := details.NewBox(&display.Screen{})
	if b.ViewPort.Height != 0 {
		t.Error("invalid object")
	}
	b.Build(entities.Task{}, false)
	if b.ViewPort.Height == 25 {
		t.Error("invalid build")
	}
}

func TestBoxUpdate(t *testing.T) {
	b := details.NewBox(&display.Screen{})
	if b.Init() != nil {
		t.Error("invalid init")
	}
	for _, key := range []tea.KeyMsg{{Type: tea.KeyUp}, {Type: tea.KeyDown}, {Runes: []rune{'G'}}, {Runes: []rune{'g'}}} {
		b.FocusIndex = 0
		for item, idx := range []int{0, 1, 2, 3} {
			if item == 0 {
				b.FocusIndex = 2
			}
			last := b.FocusIndex
			b.Blur()
			b.Update(key)
			if b.FocusIndex != last {
				t.Errorf("focus changed: %v %d% d", key, idx, b.FocusIndex)
			}
			b.Focus()
			b.FocusIndex = idx
			b.Update(key)
			if b.FocusIndex == last {
				t.Errorf("focus unchanged: %v %d% d", key, idx, b.FocusIndex)
			}
		}
	}
}

func TestDetails(t *testing.T) {
	b := details.NewBox(&display.Screen{})
	if b.Init() != nil {
		t.Error("invalid init")
	}
	v := b.View()
	if strings.TrimSpace(v) == "" {
		t.Errorf("invalid view: %s", v)
	}
	b.Focus()
	if !b.Focused() {
		t.Error("invalid focus")
	}
	b.Blur()
	if b.Focused() {
		t.Error("invalid focus")
	}
	b.Next()
	if b.FocusIndex != 1 {
		t.Errorf("invalid focus: %d", b.FocusIndex)
	}
	b.End()
	if b.FocusIndex != 3 {
		t.Errorf("invalid focus: %d", b.FocusIndex)
	}
	b.Previous()
	if b.FocusIndex != 2 {
		t.Errorf("invalid focus: %d", b.FocusIndex)
	}
	b.Start()
	if b.FocusIndex != 0 {
		t.Errorf("invalid focus: %d", b.FocusIndex)
	}
}
