// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"testing"

	"strings"
	"time"
)

func TestDateFormat_正常に取得できる(t *testing.T) {
	var verify = "2015-04-08 11:35:43.974"
	s := DateFormat(time.Now())
	if len(s) != len(verify) {
		t.Errorf("フォーマット不正[%s]、正しいフォーマットは[%s]です。", s)
	} else if !strings.Contains(s, "-") {
		t.Errorf("フォーマット不正[%s]、正しいフォーマットは[%s]です。", s)
	} else if !strings.Contains(s, ":") {
		t.Errorf("フォーマット不正[%s]、正しいフォーマットは[%s]です。", s)
	}
}

func TestDateJoblogFormat_正常に取得できる(t *testing.T) {
	var verify = "yyyyMMddHHmmss.sss"
	s := DateJoblogFormat(time.Now())
	if len(s) != len(verify) {
		t.Errorf("フォーマット不正[%s]、正しいフォーマットは[%s]です。", s, verify)
	}
}

func TestCreateFromDate_From日付が取得できる(t *testing.T) {

	s := CreateFromDate("20150414")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "2015-04-14 00:00:00.000" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "2015-04-14 00:00:00.000", s)
	}
	s = CreateFromDate("2015041A")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateFromDate("2015041")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateFromDate("X0150410")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateFromDate("20150414１")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateFromDate("2015Q414")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateFromDate("00000000")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "0000-00-00 00:00:00.000" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "0000-00-00 00:00:00.000", s)
	}
	s = CreateFromDate("99999999")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "9999-99-99 00:00:00.000" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "9999-99-99 00:00:00.000", s)
	}
}

func TestCreateToDate_To日付が取得できる(t *testing.T) {

	s := CreateToDate("20150414")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "2015-04-14 99:99:99.999" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "2015-04-14 99:99:99.999", s)
	}
	s = CreateToDate("2015041A")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateToDate("2015041")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateToDate("X0150410")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateToDate("20150414１")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateToDate("2015Q414")
	if len(s) != 0 {
		t.Errorf("空文字列が返ってくるべきところ、[%s]が返ってきました。", s)
	}
	s = CreateToDate("00000000")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "0000-00-00 99:99:99.999" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "0000-00-00 99:99:99.999", s)
	}
	s = CreateToDate("99999999")
	if len(s) == 0 {
		t.Error("空文字列が返ってきました。")
	}
	if s != "9999-99-99 99:99:99.999" {
		t.Errorf("[%v]が返ってくるべきところ、[%v]が返ってきました。", "9999-99-99 99:99:99.999", s)
	}
}
