package testutil

import (
	"fmt"
	"os"
	"testing"
)

func TestCapturer_標準出力をキャプチャできる(t *testing.T) {
	c := NewStdoutCapturer()
	c.Start()
	fmt.Println("test")
	output := c.Stop()
	if output != "test\n" {
		t.Errorf("キャプチャ結果[%s]が想定と違う。", output)
	}
}

func TestCapturer_標準エラー出力をキャプチャできる(t *testing.T) {
	c := NewStderrCapturer()
	c.Start()
	fmt.Fprintln(os.Stderr, "test")
	output := c.Stop()
	if output != "test\n" {
		t.Errorf("キャプチャ結果[%s]が想定と違う。", output)
	}
}
