package inputs_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/backend"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/inputs"
	"github.com/enckse/mayhem/internal/tui/messages"
)

type mockDB struct{}

func (m *mockDB) Add(string, any) {
}

func (m *mockDB) AddChild(string, string, any) {
}

func (m *mockDB) Remove(string) {
}

func (m *mockDB) RemoveChild(string, string) {
}

func (m *mockDB) Get() []backend.Data {
	return nil
}

func (m *mockDB) Errors() []string {
	return nil
}

func (m *mockDB) Log(_ string, _ error) {
}

func TestStackForm(t *testing.T) {
	ctx := &state.Context{}
	ctx.DB = &mockDB{}
	s := inputs.NewStackForm(entities.Stack{}, ctx)
	s.HelpKeys()
	if s.Init() != nil {
		t.Error("invalid init")
	}
	v := s.View()
	if !strings.Contains(v, "Stack Title ") {
		t.Errorf("invalid view: %s", v)
	}
	_, m := s.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if _, ok := m().(messages.Main); !ok {
		t.Error("invalid form")
	}
	_, m = s.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if m == nil {
		t.Error("invalid result")
	}
	_, m = s.Update("")
	if m != nil {
		t.Error("invalid result")
	}
	_, m = s.Update(messages.Form{Value: "xxx"})
	val, ok := m().(messages.Main)
	if !ok {
		t.Error("invalid form")
	}
	if val.Value != "refresh" {
		t.Errorf("invalid form handle: %v", val)
	}
}

func TestTaskForm(t *testing.T) {
	ctx := &state.Context{}
	ctx.DB = &mockDB{}
	ctx.Screen = &display.Screen{}
	idx := 0
	for idx < 100 {
		var now time.Time
		if idx > 50 {
			now = time.Now()
		}
		ofType := idx % 10
		var send any
		switch ofType {
		case 0, 1:
			send = "xyz"
		case 2:
			send = definitions.KeyValue{}
		case 3:
			send = time.Now()
		}

		s := inputs.NewTaskForm(entities.Task{Deadline: now}, ofType, ctx)
		idx++
		_, m := s.Update(messages.Form{Value: send})
		val, ok := m().(messages.Main)
		if !ok {
			t.Error("invalid form")
		}
		if val.Value != "refresh" {
			t.Errorf("invalid form handle: %v", val)
		}
	}
	s := inputs.NewTaskForm(entities.Task{}, 0, ctx)
	s.HelpKeys()
	if s.Init() != nil {
		t.Error("invalid init")
	}
	v := s.View()
	if !strings.Contains(v, "Task Title ") {
		t.Errorf("invalid view: %s", v)
	}
	_, m := s.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if _, ok := m().(messages.Main); !ok {
		t.Error("invalid form")
	}
	_, m = s.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if m == nil {
		t.Error("invalid result")
	}
	_, m = s.Update("")
	if m != nil {
		t.Error("invalid result")
	}
}
