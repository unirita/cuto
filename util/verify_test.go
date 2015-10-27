// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"testing"
)

func TestJobnameHasInvalidRune_禁則文字を含まない(t *testing.T) {
	target := "ソ !#%'()-=^~@`[{]};+,. _"
	if JobnameHasInvalidRune(target) {
		t.Error("禁則文字を含んでいないが失敗した。")
	}
}

func TestJobnameHasInvalidRune_バックスラッシュを含む(t *testing.T) {
	target := "abc\\"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_スラッシュを含む(t *testing.T) {
	target := "///////////////////////////////////////////////"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_コロンを含む(t *testing.T) {
	target := "a:bc"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_アスタリスクを含む(t *testing.T) {
	target := "*abc"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_クエスチョンマークを含む(t *testing.T) {
	target := "abc?"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_二重引用符を含む(t *testing.T) {
	target := "ab\"c"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_大なりを含む(t *testing.T) {
	target := "abc>"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_小なりを含む(t *testing.T) {
	target := "<abc"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_ドルマークを含む(t *testing.T) {
	target := "$abc"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}

func TestJobnameHasInvalidRune_アンパサンドを含む(t *testing.T) {
	target := "&abc"
	if !JobnameHasInvalidRune(target) {
		t.Errorf("禁則文字[%v]を含んでいるが、trueが返った。", target)
	}
}
