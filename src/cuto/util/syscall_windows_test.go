// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import "testing"

func TestLoadDll_DLLロード成功(t *testing.T) {
	dll := loadDLL("kernel32.dll")
	if dll == nil {
		t.Error("kernel32のロードに失敗しました。")
	}
	recover()
}

func TestFindProc_WinAPIを指定(t *testing.T) {
	dll := loadDLL("kernel32.dll")
	if dll == nil {
		t.Fatal("DLLロード失敗。")
	}
	p := dll.findProc("GetVersion")
	if p == nil {
		t.Error("GetVersionInfo()のローディングに失敗しました。")
	}
	recover()
}

//func TestLoadDll_DLLロード失敗(t *testing.T) {
//	dll := loadDLL("xxx.dll")
//	if dll != nil {
//		t.Error("成功するはずのないファイルがロードされました。")
//	}
//	recover()
//}

//func TestFindProc_存在しないWinAPIを指定(t *testing.T) {
//	dll := loadDLL("kernel32.dll")
//	if dll == nil {
//		t.Fatal("DLLロード失敗。")
//	}
//	p := dll.findProc("a")
//	if p != nil {
//		t.Error("成功するはずのない処理が成功しました。")
//	}
//	recover()
//}
