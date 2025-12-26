package entities_test

import (
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/entities"
)

func TestSaveDeleteTask(t *testing.T) {
	m := &mockDB{}
	s := entities.Task{}
	s.Save(m)
	if m.last == nil {
		t.Error("no save")
	}
	s.Delete(m)
	if m.last != nil {
		t.Error("no delete")
	}
}

func TestSortTasks(t *testing.T) {
	s := []entities.Task{{Title: "X00", IsFinished: true, Deadline: time.Now()}, {Title: "X01", IsFinished: false, Deadline: time.Now()}, {Title: "X10", IsFinished: true}, {Title: "X11", IsFinished: false}}
	entities.SortTasks(s)
	if len(s) != 4 || s[0].Title != "X01" || s[1].Title != "X11" || s[2].Title != "X00" || s[3].Title != "X10" {
		t.Errorf("invalid sort: %v", s)
	}
}
