// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package job

import (
	"fmt"
	"regexp"
	"strings"
)

// 引数のセパレータ用正規表現オブジェクト
var paramRegex = readyRegex()

// 正規表現オブジェクトの生成
// 引数は半角スペースで分割するが、二重引用符で括られている半角スペースは分割に使用しない。
// 例： xxx.exe A "B C D" E F" G"
//   arg[1] = "A"
//   arg[2] = "B C D"
//   arg[3] = "E"
//   arg[4] = "F\" G\""
func readyRegex() *regexp.Regexp {
	return regexp.MustCompile("((\"[^\"]*\")|[^ ])((\"[^\"]*\")*[^ ]?)*")
}

// ジョブの実行時引数を分割する。
// 両端が二重引用符の場合は、引用符を除外する。
func paramSplit(params string) []string {
	var p []string
	b := paramRegex.FindAll([]byte(params), -1)
	for i := 0; i < len(b); i++ {
		tmp := fmt.Sprintf("%s", b[i])
		if tmp[0] == '"' && tmp[len(tmp)-1] == '"' {
			tmp = tmp[1 : len(tmp)-1]
		}
		p = append(p, tmp)
	}
	return p
}

// 実行ファイル内に半角スペースが存在する場合に、二重引用符で括る
func shellFormat(shell string) string {
	rc := shell
	if sep := strings.IndexRune(shell, ' '); sep != -1 {
		rc = "\"" + shell + "\""
	}
	return rc
}
