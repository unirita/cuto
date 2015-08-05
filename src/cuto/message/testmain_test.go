package message

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	time.Local = time.FixedZone("JST", 9*60*60)
	os.Exit(m.Run())
}
