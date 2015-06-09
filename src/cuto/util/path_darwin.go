// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"os"
)

var rootPath = getCutoRoot()

func getCutoRoot() string {
	return os.Getenv("CUTOROOT")
}

// Rootフォルダを取得する
func GetRootPath() string {
	return rootPath
}

// 現在のフォルダパスを返す。
func GetCurrentPath() string {
	dir, _ := os.Getwd()
	return dir
}
