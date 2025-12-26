package backend_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

func TestAsJSON(t *testing.T) {
	testJSON(t, false)
	testJSON(t, true)
}

func testJSON(t *testing.T, pretty bool) {
	path := "testdata"
	os.MkdirAll(path, os.ModePerm)
	path = filepath.Join(path, "data.json")
	m := backend.NewMemoryBased(path, pretty, 0)
	m.Add("1", nil)
	m.Add("2", 1)
	m.AddChild("1", "2", nil)
	m.AddChild("2", "x", 5)
	m.AddChild("2", "2", 5)
	m.AddChild("1", "3", 5)
	m.AddChild("1", "4", 5)
	m.AddChild("1", "4", 6)
	m.RemoveChild("1", "3")
	if len(m.Errors()) != 0 {
		t.Error("invalid op")
	}
	b, _ := os.ReadFile(path)
	s := strings.TrimSpace(string(b))
	if pretty {
		parts := strings.Split(s, "\n")
		if parts[0] != "{" {
			t.Errorf("invalid pretty output: %s", parts[0])
		}
	} else {
		if s != `{"1":{"Node":null,"Children":{"4":{"Node":6,"Children":null}}},"2":{"Node":1,"Children":{"2":{"Node":5,"Children":null},"x":{"Node":5,"Children":null}}}}` {
			t.Error("invalid output")
		}
	}
}

func TestLoad(t *testing.T) {
	dir := "testdata"
	type parent *int
	type child *int
	const data = `{"1":{"Node":null,"Children":{"4":{"Node":6,"Children":null}}},"2":{"Node":1,"Children":{"2":{"Node":5,"Children":null},"x":{"Node":5,"Children":null}}}}`
	os.MkdirAll(dir, os.ModePerm)
	path := filepath.Join(dir, "load.invalid.json")
	m := backend.NewMemoryBased(path, false, 0)
	if err := backend.Load[parent, child](m); err == nil {
		t.Error("invalid load")
	}
	path = filepath.Join(dir, "load.json")
	os.WriteFile(path, []byte(data), 0o644)
	m = backend.NewMemoryBased(path, false, 0)
	if err := backend.Load[parent, child](m); err != nil {
		t.Error("invalid load")
	}
	if len(m.Get()) != 2 {
		t.Error("invalid get")
	}
	m.AddChild("2", "4", 6)
	b, _ := os.ReadFile(path)
	s := strings.TrimSpace(string(b))
	// Make sure MOVE still works
	if s != `{"1":{"Node":null,"Children":{}},"2":{"Node":1,"Children":{"2":{"Node":5,"Children":null},"4":{"Node":6,"Children":null},"x":{"Node":5,"Children":null}}}}` {
		t.Errorf("invalid output: %s", s)
	}
}
