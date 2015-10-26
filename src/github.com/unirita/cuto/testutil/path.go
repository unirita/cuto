package testutil

import (
	"os"
	"path/filepath"
)

var baseDir = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "unirita", "cuto")

func GetBaseDir() string {
	return baseDir
}
