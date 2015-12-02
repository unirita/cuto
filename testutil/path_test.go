package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBaseDir(t *testing.T) {
	expected := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "unirita", "cuto")
	if GetBaseDir() != expected {
		t.Errorf("GetBaseDir() => %s, want %s", GetBaseDir(), expected)
	}
}
