// Package backend provides storage solutions
package backend

type (
	// Map is a key/value (id->data) map
	Map map[string]Data
	// Data is the actual payload for data
	Data struct {
		Node     any
		Children Map
	}
	// Store defines the interface for interacting with a backend
	Store interface {
		Add(string, any)
		AddChild(string, string, any)
		Remove(string)
		RemoveChild(string, string)
		Get() []Data
		Errored() bool
		Log(string, error)
	}
)
