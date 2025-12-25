package timepicker_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/tui/inputs/timepicker"
	"github.com/enckse/mayhem/internal/tui/messages"
)

func TestInput(t *testing.T) {
	obj := timepicker.New(time.Now())
	if obj.Init() != nil {
		t.Error("invalid result")
	}
	for k, v := range []tea.KeyMsg{
		{Type: tea.KeyRight},
		{Type: tea.KeyLeft},
	} {
		val := (k + 1) + (int(time.Now().Unix()) % 100)
		for val > 0 {
			obj.Update(v)
			for _, m := range []tea.KeyMsg{{Runes: []rune{'k'}}, {Runes: []rune{'j'}}} {
				upDown := (int(time.Now().Unix()) % 100)
				for upDown > 0 {
					obj.Update(m)
					upDown--
				}
			}
			val--
		}
	}
	_, val := obj.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_, ok := val().(messages.Form)
	if !ok {
		t.Error("invalid results")
	}
	v := obj.View()
	if !strings.Contains(v, "YYYY") {
		t.Error("invalid view")
	}
}
