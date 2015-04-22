// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"fmt"
	"os"
	"strings"
	"unsafe"
)

var (
	modulePath = getModulePath()
)

const max_path = 520

// Rootフォルダを取得する（実行ファイルはRootパス下のbinフォルダにあると想定）
func GetRootPath() string {
	c := modulePath[:strings.LastIndex(modulePath, fmt.Sprintf("%c", os.PathSeparator))]
	rootPath := modulePath[:strings.LastIndex(c, fmt.Sprintf("%c", os.PathSeparator))]
	//	fmt.Println("Debug : rootPath = ", rootPath)
	return rootPath
}

// 現在のフォルダパスを返す。
func GetCurrentPath() string {
	c := modulePath[:strings.LastIndex(modulePath, fmt.Sprintf("%c", os.PathSeparator))]
	//	fmt.Println("Debug : currentPath = ", c)
	return c
}

func getModulePath() string {
	// MAX_PATHがUTF-8になる場合は、これくらいあれば十分か？
	var buf [max_path]byte
	procGetModuleFileNameW.Call(0, uintptr(unsafe.Pointer(&buf)), (uintptr)(max_path))

	// Unicodeで取得しているので、2byte目が0の部分を除外する。
	var path [max_path / 2]byte
	var j int
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			path[j] = buf[i]
			j++
		}
	}
	return fmt.Sprintf("%s", path)
}
