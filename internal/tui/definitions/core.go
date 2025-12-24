// Package definitions defines common types/utilities
package definitions

// KeyValue is a simple key/value object
type KeyValue struct {
	Key   uint
	Value string
}

const (
	// StackTitleIndex is the title index for stack fields (indexed)
	StackTitleIndex int = iota
)

const (
	// TaskTitleIndex is the title index for task fields (indexed)
	TaskTitleIndex int = iota
	// TaskNotesIndex is the notes index for task fields (indexed)
	TaskNotesIndex
	// TaskPriorityIndex is the priority index for task fields (indexed)
	TaskPriorityIndex
	// TaskDeadlineIndex is the deadline index for task fields (indexed)
	TaskDeadlineIndex
)

const (
	// TaskLastIndex is the last known task item (index)
	TaskLastIndex = TaskDeadlineIndex
	// IsDelete is a delete command
	IsDelete = "delete"
	// IsMove is a move command
	IsMove = "move"
)
