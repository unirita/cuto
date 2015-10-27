// 入力チェック。
// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"strings"
)

const invalid_jobname_ptr = "\\/:*?\"<>|$&"

// 指定したジョブ名に、禁止文字が存在するか確認する。
// 現在の仕様では 「 \/:*?"<>|$ 」の記号が使用禁止。
//
// param : jobname ジョブ名。
//
// return : 禁止文字が含まれている場合はtrueを返す。
func JobnameHasInvalidRune(jobname string) bool {
	return -1 != strings.IndexAny(jobname, invalid_jobname_ptr)
}

// 他に入力項目の禁則文字チェックを行いたい場合は、このファイルへ追加する。
