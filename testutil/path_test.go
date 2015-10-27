package testutil

import (
	"os"
	"testing"
)

func TestBaseDir(t *testing.T) {
	expected := os.Getenv("GOPATH")
	if os.PathSeparator == '\\' {
		expected += `\src\github.com\unirita\cuto`
	} else {
		expected += `/src/github.com/unirita/cuto`
	}

	if GetBaseDir() != expected {
		t.Errorf("GetBaseDir() => %s, want %s", GetBaseDir(), expected)
	}
}
