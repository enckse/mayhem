package entities_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/backend"
	"github.com/enckse/mayhem/internal/entities"
)

type mockDB struct {
	last any
	data []backend.Data
	cat  string
}

func (m *mockDB) Add(_ string, obj any) {
	m.last = obj
}

func (m *mockDB) AddChild(_, _ string, obj any) {
	m.last = obj
}

func (m *mockDB) Remove(string) {
	m.last = nil
}

func (m *mockDB) RemoveChild(string, string) {
	m.last = nil
}

func (m *mockDB) Get() []backend.Data {
	return m.data
}

func (m *mockDB) Errors() []string {
	return nil
}

func (m *mockDB) Log(cat string, _ error) {
	m.cat = cat
}

func TestFindByIndex(t *testing.T) {
	idx := entities.FindByIndex([]entities.Stack{}, "")
	if idx != -1 {
		t.Errorf("invalid index: %d", idx)
	}
	idx = entities.FindByIndex([]entities.Stack{{ID: "id"}}, "id")
	if idx != 0 {
		t.Errorf("invalid index: %d", idx)
	}
	idx = entities.FindByIndex([]entities.Stack{{}}, "1")
	if idx != -1 {
		t.Errorf("invalid index: %d", idx)
	}
}
