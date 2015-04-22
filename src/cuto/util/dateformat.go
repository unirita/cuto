// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"fmt"
	"regexp"
	"time"
)

// 日時を「yyyy-MM-dd HH:mm:ss.sss」に変換する。
func DateFormat(t time.Time) string {
	milisec := fmt.Sprintf("%09d", t.Nanosecond())[:3]
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), milisec)
}

// 日時を「yyyyMMddHHmmss.sss」に変換する。
func DateJoblogFormat(t time.Time) string {
	milisec := fmt.Sprintf("%09d", t.Nanosecond())[:3]
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d.%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), milisec)
}

// FROM日付を取得する。日付としての正誤は確認しない。
// 引数に指定した値が8文字に満たない、または数値以外が含まれている場合は空文字列を返す。
func CreateFromDate(from string) string {
	if len(from) != 8 {
		return ""
	}
	if isMatch, _ := regexp.MatchString("^[0-9]{8}$", from); !isMatch {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s 00:00:00.000", from[0:4], from[4:6], from[6:8])
}

// TO日付を取得する。日付としての正誤は確認しない。
// 引数に指定した値が8文字に満たない、または数値以外が含まれている場合は空文字列を返す。
func CreateToDate(to string) string {
	if len(to) != 8 {
		return ""
	}
	if isMatch, _ := regexp.MatchString("^[0-9]{8}$", to); !isMatch {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s 99:99:99.999", to[0:4], to[4:6], to[6:8])
}
