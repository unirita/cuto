// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"os"
	"testing"
)

func TestGetRootPath_正常にRootPathが取得できる(t *testing.T) {
	os.Setenv("CUTOROOT", "UNITTest")
	s := GetRootPath()
	if len(s) == 0 {
		t.Error("取得失敗。")
	}
}

func TestGetCurrentPath_正常にCurrentPathが取得できる(t *testing.T) {
	s := GetCurrentPath()
	if len(s) == 0 {
		t.Error("取得失敗。")
	}
}
