package job

import (
	"testing"

	"github.com/unirita/cuto/testutil"
)

func TestStdoutBufferWrite_OutputEnabled(t *testing.T) {
	buf := NewStdoutBuffer(true)

	c := testutil.NewStdoutCapturer()
	c.Start()
	buf.Write([]byte("test1"))
	buf.Write([]byte("test2"))
	out := c.Stop()

	written := buf.String()
	if written != "test1test2" {
		t.Errorf("buf.String() => %s, wants %s", written, "test1test2")
	}
	if out != written {
		t.Errorf("Buffered string did not output for stdout.")
		t.Logf("Output: %s", out)
	}
}

func TestStdoutBufferWrite_OutputDisabled(t *testing.T) {
	buf := NewStdoutBuffer(false)

	c := testutil.NewStdoutCapturer()
	c.Start()
	buf.Write([]byte("test1"))
	buf.Write([]byte("test2"))
	out := c.Stop()

	written := buf.String()
	if written != "test1test2" {
		t.Errorf("buf.String() => %s, wants %s", written, "test1test2")
	}
	if out != "" {
		t.Errorf("Unexpected string was output for stdout.")
		t.Logf("Output: %s", out)
	}
}
