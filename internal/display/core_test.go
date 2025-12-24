package display_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/enckse/mayhem/internal/display"
)

func TestNewScreen(t *testing.T) {
	s := display.NewScreen()
	if s.Table.ViewHeight != 25 {
		t.Error("invalid view height")
	}
}

func TestInputFormStyle(t *testing.T) {
	s := display.NewScreen()
	s.Width = 50
	style := s.InputFormStyle()
	if style.GetWidth() != 48 {
		t.Error("invalid style result")
	}
}

func TestTableStyle(t *testing.T) {
	style := display.TableStyle(display.StackTableType)
	if fmt.Sprintf("%v", style.Selected.GetBackground()) != "#019187" {
		t.Errorf("invalid style: %v", style.Selected.GetBackground())
	}
	style = display.TableStyle(display.TaskTableType)
	if fmt.Sprintf("%v", style.Selected.GetBackground()) != "#f1b44c" {
		t.Errorf("invalid style: %v", style.Selected.GetBackground())
	}
}

func TestEmptyTaskStyle(t *testing.T) {
	s := display.NewScreen()
	style := s.EmptyTaskStyle()
	if style.GetWidth() != display.TaskTableWidth || style.GetHeight() != 26 {
		t.Error("invalid style")
	}
}

func TestEmptyDetailsStyle(t *testing.T) {
	s := display.NewScreen()
	style := s.EmptyDetailsStyle()
	if style.GetHeight() != 26 {
		t.Error("invalid style")
	}
}

func TestDetailsBoxWidth(t *testing.T) {
	s := display.NewScreen()
	style := s.DetailsBoxWidth()
	if style != -93 {
		t.Errorf("invalid value %d", style)
	}
}

func TestDetailsBoxHeight(t *testing.T) {
	s := display.NewScreen()
	style := s.DetailsBoxHeight()
	if style != 25 {
		t.Errorf("invalid value %d", style)
	}
}

func TestDetailsBoxStyle(t *testing.T) {
	s := display.NewScreen()
	s.Width = 100
	style := s.DetailsBoxStyle()
	if style.GetWidth() != 7 || style.GetHeight() != 25 {
		t.Errorf("invalid style result %d %d", style.GetWidth(), style.GetHeight())
	}
}

func TestDetailsItemStyle(t *testing.T) {
	s := display.NewScreen()
	s.Width = 100
	style := s.DetailsItemStyle(false)
	if style.GetWidth() != 5 {
		t.Errorf("invalid style result %d", style.GetWidth())
	}
	if fmt.Sprintf("%v", style.GetBackground()) != "{}" {
		t.Errorf("invalid background: %s", style.GetBackground())
	}
	style = s.DetailsItemStyle(true)
	if fmt.Sprintf("%v", style.GetBackground()) != "#333c4d" {
		t.Errorf("invalid background: %s", style.GetBackground())
	}
}

func TestItemContainerStyle(t *testing.T) {
	s := display.NewScreen()
	s.Width = 100
	style := s.ItemContainerStyle(false)
	if style.GetWidth() != 7 {
		t.Errorf("invalid style result %d", style.GetWidth())
	}
	if fmt.Sprintf("%v", style.GetBackground()) != "{}" {
		t.Errorf("invalid background: %s", style.GetBackground())
	}
	style = s.ItemContainerStyle(true)
	if fmt.Sprintf("%v", style.GetBackground()) != "#333c4d" {
		t.Errorf("invalid background: %s", style.GetBackground())
	}
}

func TestEmptyTaskView(t *testing.T) {
	s := display.NewScreen()
	val := s.EmptyTaskView()
	if !strings.Contains(val, "key to explore this stack") {
		t.Errorf("invalid render: %s", val)
	}
}

func TestEmptyDetailsView(t *testing.T) {
	s := display.NewScreen()
	val := s.EmptyDetailsView()
	if !strings.Contains(val, "key to see task details") {
		t.Errorf("invalid render: %s", val)
	}
}
