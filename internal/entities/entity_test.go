package entities_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

type mockDB struct {
	last any
	ctx  *state.Context
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

func (m *mockDB) SyncJSON(ctx *state.Context) {
	m.ctx = ctx
}

func (m *mockDB) Errors() []string {
	return nil
}

func TestSync(t *testing.T) {
	ctx := &state.Context{}
	m := &mockDB{}
	ctx.DB = m
	entities.Sync(ctx, entities.Stack{})
	if m.last == nil || m.ctx != nil {
		t.Error("invalid save")
	}
	ctx.Config.JSON.Sync = true
	entities.Sync(ctx, entities.Stack{})
	if m.ctx == nil {
		t.Error("invalid save")
	}
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
