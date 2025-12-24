package keys_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/tui/keys"
)

func TestShortHelp(t *testing.T) {
	h := keys.Map{}.ShortHelp()
	if len(h) == 0 {
		t.Error("invalid short help")
	}
}

func TestFullHelp(t *testing.T) {
	h := keys.Map{}.FullHelp()
	if len(h) == 0 || len(h[0]) == 0 {
		t.Error("invalid full help")
	}
}
