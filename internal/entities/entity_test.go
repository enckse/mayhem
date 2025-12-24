package entities_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

type mockDB struct {
	last any
	err  error
	ctx  *state.Context
	id uint
}

func (m *mockDB) Save(obj any) {
	m.last = obj
}

func (m *mockDB) Create(obj any) error {
	m.last = obj
	return m.err
}

func (m *mockDB) Delete(_ any) {
	m.last = nil
}

func (m *mockDB) Find(obj any, id uint) {
	m.last = obj
	m.id = id
}

func (m *mockDB) Stacks(_ any) error {
	return m.err
}

func (m *mockDB) SyncJSON(ctx *state.Context) {
	m.ctx = ctx
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
