package message

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
