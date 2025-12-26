package backend

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// MemoryBased is a memory-based backend (that can sync to JSON file)
type MemoryBased struct {
	data      Map
	children  map[string]string
	errors    []string
	maxErrors int
	pretty    bool
	file      string
}

// NewMemoryBased will create a new memory-based backend
func NewMemoryBased(file string, pretty bool, maxErrors uint) *MemoryBased {
	return &MemoryBased{
		data:      make(Map),
		maxErrors: int(maxErrors),
		pretty:    pretty,
		file:      file,
		children:  make(map[string]string),
	}
}

// Errors will get the list of errors
func (m *MemoryBased) Errors() []string {
	return m.errors
}

// Add will add a new entity
func (m *MemoryBased) Add(id string, data any) {
	err := func() error {
		if strings.TrimSpace(id) == "" {
			return errors.New("id is empty")
		}
		v, ok := m.data[id]
		if !ok {
			v = Data{}
			v.Children = make(Map)
		}
		v.Node = data
		m.data[id] = v
		return nil
	}()
	m.Log("add", err)
}

// AddChild will add a child entity
func (m *MemoryBased) AddChild(parent, id string, data any) {
	err := func() error {
		if strings.TrimSpace(parent) == "" {
			return errors.New("parent is empty")
		}
		if strings.TrimSpace(id) == "" {
			return errors.New("id is empty")
		}
		if _, ok := m.data[parent]; !ok {
			return fmt.Errorf("parent not found: %v", parent)
		}
		if v, ok := m.children[id]; ok {
			delete(m.data[v].Children, id)
		}
		c, ok := m.data[parent].Children[id]
		if !ok {
			c = Data{}
		}
		c.Node = data
		m.data[parent].Children[id] = c
		m.children[id] = parent
		return nil
	}()
	m.Log("addchild", err)
}

// Remove will remove an entity
func (m *MemoryBased) Remove(id string) {
	err := func() error {
		delete(m.data, id)
		return nil
	}()
	m.Log("remove", err)
}

// RemoveChild will remove a child entity
func (m *MemoryBased) RemoveChild(parent, id string) {
	err := func() error {
		if _, ok := m.data[parent]; !ok {
			return fmt.Errorf("parent not found: %v", parent)
		}
		delete(m.data[parent].Children, id)
		return nil
	}()
	m.Log("removechild", err)
}

// Get will return the backing data
func (m *MemoryBased) Get() []Data {
	var data []Data
	for _, v := range m.data {
		data = append(data, v)
	}
	return data
}

// Log will add an error to the backend data
func (m *MemoryBased) Log(cat string, err error) {
	if err == nil {
		return
	}
	m.errors = append(m.errors, fmt.Sprintf("[%s] %s: %v", time.Now().Format("2006-01-02T15:04:05"), cat, err))
	if m.maxErrors > 0 && len(m.errors) > m.maxErrors {
		m.errors = m.errors[1:]
	}
}
