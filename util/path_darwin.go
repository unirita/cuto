// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"fmt"
	"os"
)

var rootPath = getCutoRoot()

func getCutoRoot() string {
	d := os.Getenv("CUTOROOT")
	if len(d) == 0 {
		panic("Not setting environment argument $CUTOROOT")
	}
	return d
}

// Rootフォルダを取得する
func GetRootPath() string {
	return rootPath
}

// 現在のフォルダパスを返す。
func GetCurrentPath() string {
	return fmt.Sprintf("%s%c%s", rootPath, os.PathSeparator, "bin")
}
