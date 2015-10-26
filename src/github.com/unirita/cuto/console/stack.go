// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package console

import (
	"github.com/unirita/cuto/log"
	"runtime"
)

func PrintStack() {
	for i := 2; ; i++ { // ここと、呼び出し元のconsoleは出さない。
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		log.Error("\t", file, ":", line)
	}
}
