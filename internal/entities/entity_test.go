package entities_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/entities"
)

type mockDB struct {
	last any
}

func (m *mockDB) Save(obj any) {
	m.last = obj
}

func (m *mockDB) Create(obj any) {
	m.last = obj
}

func (m *mockDB) Delete(_ any) {
	m.last = nil
}

func (m *mockDB) Fetch() any {
	return m.last
}

func (m *mockDB) Errors() []string {
	return nil
}

func TestFindByIndex(t *testing.T) {
	idx := entities.FindByIndex([]entities.Stack{}, 0)
	if idx != -1 {
		t.Errorf("invalid index: %d", idx)
	}
	idx = entities.FindByIndex([]entities.Stack{{}}, 0)
	if idx != 0 {
		t.Errorf("invalid index: %d", idx)
	}
	idx = entities.FindByIndex([]entities.Stack{{}}, 1)
	if idx != -1 {
		t.Errorf("invalid index: %d", idx)
	}
}

func TestEntityID(t *testing.T) {
	e := entities.Stack{}
	e.ID = 1
	if e.EntityID() != 1 {
		t.Error("invalid id")
	}
}
