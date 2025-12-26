package backend_test

import (
	"errors"
	"testing"

	"github.com/enckse/mayhem/internal/backend"
)

func TestLogErrors(t *testing.T) {
	m := backend.NewMemoryBased("", false, 0)
	m.Log("abc", nil)
	m.Log("xyz", errors.New("test"))
	m.Log("xyz", errors.New("test"))
	m.Log("xyz", errors.New("test"))
	if len(m.Errors()) != 3 {
		t.Error("invalid error set")
	}
	m = backend.NewMemoryBased("", false, 1)
	m.Log("abc", nil)
	m.Log("xyz", errors.New("test"))
	m.Log("xyz", errors.New("test"))
	m.Log("xyz", errors.New("test"))
	if len(m.Errors()) != 1 {
		t.Error("invalid error set")
	}
}

func TestAdd(t *testing.T) {
	m := backend.NewMemoryBased("", false, 0)
	m.Add("", nil)
	if len(m.Errors()) != 1 {
		t.Error("invalid add")
	}
	m.Add("x", 1)
	if len(m.Errors()) != 1 {
		t.Error("invalid add")
	}
	if len(m.Get()) != 1 {
		t.Error("invalid add")
	}
	m.Add("x", 1)
	if len(m.Errors()) != 1 {
		t.Error("invalid add")
	}
	if len(m.Get()) != 1 {
		t.Error("invalid add")
	}
	m.Add("y", 1)
	if len(m.Errors()) != 1 {
		t.Error("invalid add")
	}
	if len(m.Get()) != 2 {
		t.Error("invalid add")
	}
}

func TestAddChild(t *testing.T) {
	m := backend.NewMemoryBased("", false, 0)
	m.Add("1", 1)
	m.Add("2", 2)
	m.AddChild("", "", 1)
	if len(m.Errors()) != 1 {
		t.Error("invalid add")
	}
	m.AddChild("x", "", 1)
	if len(m.Errors()) != 2 {
		t.Error("invalid add")
	}
	m.AddChild("x", "y", 1)
	if len(m.Errors()) != 3 {
		t.Error("invalid add")
	}
	m.AddChild("1", "0", 1)
	if len(m.Errors()) != 3 {
		t.Error("invalid add")
	}
	data := m.Get()
	if len(data) != 2 {
		t.Error("invalid results")
	}
	for _, d := range data {
		length := 1
		if d.Node.(int) == 2 {
			length = 0
		}
		if len(d.Children) != length {
			t.Error("invalid results")
		}
	}
	m.AddChild("1", "1", 1)
	if len(m.Errors()) != 3 {
		t.Error("invalid add")
	}
	data = m.Get()
	if len(data) != 2 {
		t.Error("invalid results")
	}
	for _, d := range data {
		length := 2
		if d.Node.(int) == 2 {
			length = 0
		}
		if len(d.Children) != length {
			t.Error("invalid results")
		}
	}
	m.AddChild("2", "0", 1)
	if len(m.Errors()) != 3 {
		t.Error("invalid add")
	}
	data = m.Get()
	if len(data) != 2 {
		t.Error("invalid results")
	}
	for _, d := range data {
		if len(d.Children) != 1 {
			t.Error("invalid results")
		}
	}
}

func TestRemove(t *testing.T) {
	m := backend.NewMemoryBased("", false, 0)
	m.Add("1", nil)
	m.Remove("")
	if len(m.Errors()) != 0 {
		t.Error("invalid remove")
	}
	if len(m.Get()) != 1 {
		t.Error("invalid remove")
	}
	m.Remove("1")
	if len(m.Errors()) != 0 {
		t.Error("invalid remove")
	}
	if len(m.Get()) != 0 {
		t.Error("invalid remove")
	}
}

func TestRemoveChild(t *testing.T) {
	m := backend.NewMemoryBased("", false, 0)
	m.Add("1", nil)
	m.AddChild("1", "2", nil)
	m.RemoveChild("x", "")
	if len(m.Errors()) != 1 {
		t.Error("invalid remove")
	}
	m.RemoveChild("1", "")
	if len(m.Errors()) != 1 {
		t.Error("invalid remove")
	}
	if len(m.Get()[0].Children) != 1 {
		t.Error("invalid remove")
	}
	m.RemoveChild("1", "2")
	if len(m.Errors()) != 1 {
		t.Error("invalid remove")
	}
	if len(m.Get()[0].Children) != 0 {
		t.Error("invalid remove")
	}
}
