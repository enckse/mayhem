package messages_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/tui/messages"
)

func TestGoToMain(t *testing.T) {
	if _, ok := messages.MainGoTo().(messages.Main); !ok {
		t.Error("invalid message")
	}
	v := messages.MainGoToWith("abc")().(messages.Main)
	if v.Value != "abc" {
		t.Error("invalid message")
	}
}

func TestGoToForm(t *testing.T) {
	v := messages.FormGoToWith("abc")().(messages.Form)
	if v.Value != "abc" {
		t.Error("invalid message")
	}
}
