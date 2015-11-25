// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

// +build darwin linux

package job

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/message"
)

const testServantIni = "servant_l.ini"

func TestDoJobRequest_拡張子無しSHジョブが正常に実行できる(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       201,
		JID:       "serviceTask_001",
		Path:      "job",
		Param:     "XX 0",
		Env:       "TESTENV1=ENVENV+ENV0=AAAAA+ENV2=BBBB",
		Workspace: "",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}
	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	if len(<-stCh) == 0 {
		t.Error("ジョブ開始時間が送信されていない.")
	}
	close(stCh)

	if req.NID != res.NID {
		t.Error("NIDがリクエストとレスポンスで異なる.")
	}
	if req.JID != res.JID {
		t.Error("IDがリクエストとレスポンスで異なる.")
	}
	if res.RC != 0 {
		t.Errorf("ジョブが正常終了するはずなのに異常終了した. - %v", res.RC)
	}
	if len(res.Detail) > 0 {
		t.Error("ジョブが正常終了するはずなのに、エラーメッセージがある.", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("テストコード側の変数機能が未実装.")
	} else if res.Var != "ENVENV XX" {
		t.Errorf("変数の値が不正[%s].もしくは引数渡しに問題あり.", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間が無い.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間が無い.")
	}

	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログがない。", err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var env1, env2, env3 bool
	for scanner.Scan() {
		line := scanner.Text()
		if -1 != strings.Index(line, "ENVENV") {
			env1 = true
		} else if -1 != strings.Index(line, "AAAAA") {
			env2 = true
		} else if -1 != strings.Index(line, "BBBB") {
			env3 = true
		}
	}
	if !env1 || !env2 || !env3 {
		t.Error("環境変数が正常に渡っていない可能性があります。", env1, env2, env3)
	}
}

func TestDoJobRequest_パス指定ありCSHジョブが正常に実行できる(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       202,
		JID:       "serviceTask_002",
		Path:      "job.csh",
		Param:     "",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: testJobPath,
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Errorf("ジョブが正常終了するはずなのに異常終了した. - %v", res.RC)
	}
	if len(res.Detail) > 0 {
		t.Error("ジョブが正常終了するはずなのに、エラーメッセージがある.", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "ENVENVENV" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間が無い.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間が無い.")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログがない。", err.Error())
	}
	defer file.Close()
}

func TestDoJobRequest_存在しないジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       203,
		JID:       "serviceTask_003",
		Path:      "nothing.bat",
		Param:     "",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: testJobPath,
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("実行失敗の場合、RCは0のはず.")
	}
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	}
	if len(res.Var) > 0 {
		t.Error("テストコード側の変数機能が未実装.")
	}
	if len(res.St) != 0 {
		t.Error("ジョブ開始時間がある.")
	}
	if len(res.Et) != 0 {
		t.Error("ジョブ終了時間がある.")
	}
}
func TestDoJobRequest_RCで警告終了するジョブ_閾値と同じ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       204,
		JID:       "serviceTask_004",
		Path:      "job.bash",
		Param:     "X 4",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: testJobPath,
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 4 {
		t.Errorf("RCは4のはず. - %v", res.RC)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "ENVENVENV X" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.WARN {
		t.Error("statが警告終了ではない")
	}
	if res.Detail != detailWarnRC {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}
func TestDoJobRequest_RCで異常終了するジョブ_閾値と同じ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       205,
		JID:       "serviceTask_004",
		Path:      "job.bash",
		Param:     "X 12",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: testJobPath,
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 12 {
		t.Errorf("RCは12のはず. - %v", res.RC)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "ENVENVENV X" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.ABNORMAL {
		t.Error("statが異常終了ではない")
	}
	if res.Detail != detailErrRC {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_標準出力で警告終了するジョブ_RC確認なし(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       206,
		JID:       "serviceTask_004",
		Path:      "job.ksh",
		Param:     "\"A B\"",
		Env:       "TESTENV1=!!!WARNING!!!",
		Workspace: testJobPath,
		WarnRC:    0,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("RCは0のはず.")
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "!!!WARNING!!! A B" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.WARN {
		t.Error("statが警告終了ではない")
	}
	if res.Detail != detailWarnPtn {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}
func TestDoJobRequest_標準エラー出力で異常終了するzshジョブ_RC確認なし(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       207,
		JID:       "serviceTask_006",
		Path:      "stderr.zsh",
		Param:     "!!!ERROR",
		Env:       "",
		Workspace: testJobPath,
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("RCは0のはず.")
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "Argument1=!!!ERROR" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.ABNORMAL {
		t.Error("statが異常終了ではない")
	}
	if res.Detail != detailErrPtn {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_タイムアウト時間を超えたら異常終了する(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       208,
		JID:       "serviceTask_001",
		Path:      "twosec.sh",
		Param:     "",
		Env:       "",
		Workspace: testJobPath,
		WarnRC:    0,
		WarnStr:   "",
		ErrRC:     0,
		ErrStr:    "",
		Timeout:   1,
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Errorf("戻り値が想定外 - %v", res.RC)
	}
	if res.Stat != db.ABNORMAL {
		t.Error("異常終了するはずが、正常終了している。")
	}
	if len(res.Detail) == 0 {
		t.Error("異常メッセージがセットされていない。")
	}
}
