// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package testutil

import (
	"bytes"
	"io"
	"os"
)

// 標準出力・標準エラー出力のキャプチャを行う構造体
type Capturer struct {
	isStderr bool
	original *os.File
	bufCh    chan string
	out      *os.File
	in       *os.File
}

// 標準出力をキャプチャするCapturerを生成する。
//
// return : 生成したCapturerオブジェクト
func NewStdoutCapturer() *Capturer {
	c := new(Capturer)
	c.isStderr = false
	return c
}

// 標準エラー出力をキャプチャするCapturerを生成する。
//
// return : 生成したCapturerオブジェクト
func NewStderrCapturer() *Capturer {
	c := new(Capturer)
	c.isStderr = true
	return c
}

// キャプチャを開始する。
func (c *Capturer) Start() {
	if c.isStderr {
		c.original = os.Stderr
	} else {
		c.original = os.Stdout
	}
	var err error
	c.in, c.out, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	if c.isStderr {
		os.Stderr = c.out
	} else {
		os.Stdout = c.out
	}
	c.bufCh = make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, c.in)
		c.bufCh <- b.String()
	}()
}

// キャプチャを停止する。
//
// return : キャプチャ開始から停止までの間に出力された文字列
func (c *Capturer) Stop() string {
	c.out.Close()
	defer c.in.Close()
	if c.isStderr {
		os.Stderr = c.original
	} else {
		os.Stdout = c.original
	}
	return <-c.bufCh
}
