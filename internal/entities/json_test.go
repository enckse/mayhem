package entities_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

func TestDumpJSON(t *testing.T) {
	m := &mockDB{}
	ctx := &state.Context{}
	ctx.DB = m
	var buf bytes.Buffer
	if err := entities.DumpJSON(&buf, ctx); err != nil {
		t.Errorf("invalid error: %v", err)
	}
	if !strings.Contains(buf.String(), "New Stack") {
		t.Errorf("invalid dump: %s", buf.String())
	}
}

func TestLoadJSON(t *testing.T) {
	m := &mockDB{}
	ctx := &state.Context{}
	ctx.DB = m
	reader := strings.NewReader("[{\"Title\": \"xyz\"}]")
	if err := entities.LoadJSON(ctx, false, reader); err != nil {
		t.Errorf("invalid error: %v", err)
	}
	if !strings.Contains(fmt.Sprintf("%v", m.last), "xyz") {
		t.Error("invalid load")
	}
	reader = strings.NewReader("[{\"Title\": \"zzz\"}]")
	if err := entities.LoadJSON(ctx, true, reader); err != nil {
		t.Errorf("invalid error: %v", err)
	}
	if !strings.Contains(fmt.Sprintf("%v", m.last), "zzz") {
		t.Error("invalid load")
	}
}
