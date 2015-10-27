// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package console

import (
	"testing"

	"github.com/unirita/cuto/testutil"
)

func TestDisplay_メッセージを出力できる_引数なし(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	c.Start()

	Display("CTM003E")

	output := c.Stop()

	if output != "CTM003E INVALID ARGUMENT.\n" {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestDisplay_メッセージを出力できる_引数あり(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	c.Start()

	Display("CTM019E", "something error.")

	output := c.Stop()

	if output != "CTM019E EXCEPTION OCCURED - something error.\n" {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestDisplayError_メッセージをエラー出力できる_引数なし(t *testing.T) {
	c := testutil.NewStderrCapturer()
	c.Start()

	DisplayError("CTM003E")

	output := c.Stop()

	if output != "CTM003E INVALID ARGUMENT.\n" {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestDisplayError_メッセージをエラー出力できる_引数あり(t *testing.T) {
	c := testutil.NewStderrCapturer()
	c.Start()

	DisplayError("CTM019E", "something error.")

	output := c.Stop()

	if output != "CTM019E EXCEPTION OCCURED - something error.\n" {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestGetMessage_メッセージを文字列として取得できる_引数なし(t *testing.T) {
	msg := GetMessage("CTM003E")
	if msg != "CTM003E INVALID ARGUMENT." {
		t.Errorf("取得値[%s]が想定と違います。", msg)
	}
}

func TestGetMessage_メッセージを文字列として取得できる_引数あり(t *testing.T) {
	msg := GetMessage("CTM019E", "something error.")
	if msg != "CTM019E EXCEPTION OCCURED - something error." {
		t.Errorf("取得値[%s]が想定と違います。", msg)
	}
}
