// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package job

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"cuto/db"
	"cuto/message"

	"cuto/servant/config"
)

var conf *config.ServantConfig

// ジョブログなどの掃除
func init() {
	time.Local = time.FixedZone("JST", 9*60*60)

	s := os.PathSeparator
	testPath := fmt.Sprintf("%s%c%s%c%s%c%s%c%s", os.Getenv("GOPATH"), s, "test", s, "cuto", s, "servant", s, "job")
	err := os.Chdir(testPath)
	config.RootPath = testPath
	if err != nil {
		panic(err.Error())
	}
	configPath := fmt.Sprintf("%s%c%s", testPath, s, "servant.ini")
	conf = config.ReadConfig(configPath)
	os.RemoveAll(conf.Dir.JoblogDir)
	err = os.Mkdir(conf.Dir.JoblogDir, 0666)
	if err != nil {
		panic(err.Error())
	}
}

// メッセージ・DB向けの時刻フォーマットをジョブログファイル名向けのフォーマットに変換する。
func stToLocalTimestamp(st string) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", st, time.UTC)
	if err != nil {
		panic("Unexpected time format: " + err.Error())
	}
	return t.Local().Format("20060102150405.000")
}

// ジョブログファイル名をフルパスで作成する。
// ”開始日(YYYYMMDD)\インスタンスID.ジョブ名（拡張子なし）.開始日時（yyyyMMddHHmmss.sss).log
func createJoblogFileName(req *message.Request, st string, nID int, jID string) string {
	var job string
	if strings.LastIndex(req.Path, "\\") != -1 {
		tokens := strings.Split(req.Path, "\\")
		job = tokens[len(tokens)-1]
	} else if strings.LastIndex(req.Path, "/") != -1 {
		tokens := strings.Split(req.Path, "/")
		job = tokens[len(tokens)-1]
	} else {
		job = req.Path
	}
	if extpos := strings.LastIndex(job, "."); extpos != -1 {
		job = job[:extpos]
	}
	joblogDir := fmt.Sprintf("%v%c%v", conf.Dir.JoblogDir, os.PathSeparator, st[:8])
	return fmt.Sprintf("%v%c%v.%v.%v.%v.log", joblogDir, os.PathSeparator, nID, job, jID, st)
}

func TestDoJobRequest_ジョブが正常に実行できる(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       101,
		JID:       "serviceTask_001",
		Path:      "job.bat",
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
		t.Errorf("ジョブが正常終了するはずなのに異常終了した. - ", res.RC)
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

func TestDoJobRequest_パス指定あり引数なしジョブが正常に実行できる(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       102,
		JID:       "serviceTask_002",
		Path:      "job.bat",
		Param:     "",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}
	gopath := os.Getenv("GOPATH")
	req.Path = fmt.Sprintf("%s%c%s%c%s%c%s%c%s%c%s%c%s",
		gopath, os.PathSeparator, "test", os.PathSeparator,
		"cuto", os.PathSeparator, "servant", os.PathSeparator, "job", os.PathSeparator,
		"jobscript", os.PathSeparator, "job.bat")

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("ジョブが正常終了するはずなのに異常終了した.")
	}
	if len(res.Detail) > 0 {
		t.Error("ジョブが正常終了するはずなのに、エラーメッセージがある.", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "ENVENVENV " {
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
		NID:       103,
		JID:       "serviceTask_003",
		Path:      "nothing.bat",
		Param:     "",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
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
		NID:       104,
		JID:       "serviceTask_004",
		Path:      "job.bat",
		Param:     "X 4",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 4 {
		t.Error("RCは4のはず.")
	}
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	} else if res.Detail != detailWarnRC {
		t.Errorf("想定外のメッセージ - %s", res.Detail)
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで警告終了するがチェックしないジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       105,
		JID:       "serviceTask_004",
		Path:      "job.bat",
		Param:     "X 4",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    0,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 4 {
		t.Error("RCは4のはず.")
	}
	if len(res.Detail) != 0 {
		t.Error("異常終了メッセージが存在しない.")
	}
	if len(res.Var) == 0 {
		t.Error("テストコード側の変数機能が未実装.")
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.NORMAL {
		t.Error("statが正常終了ではない")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで警告終了しないジョブ_閾値未満(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       106,
		JID:       "serviceTask_004",
		Path:      "job.bat",
		Param:     "X 3",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 3 {
		t.Error("RCは3のはず.")
	}
	if len(res.Detail) != 0 {
		t.Error("異常終了メッセージが存在しない.")
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
	if res.Stat != db.NORMAL {
		t.Error("statが正常終了ではない")
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
		NID:       107,
		JID:       "serviceTask_004",
		Path:      "job.bat",
		Param:     "\"A B\"",
		Env:       "TESTENV1=!!!WARNING!!!",
		Workspace: "C:\\Go",
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
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	} else if res.Detail != detailWarnPtn {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "!!!WARNING!!! \"A B\"" {
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_標準出力で警告終了するがチェックしないジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       108,
		JID:       "serviceTask_004",
		Path:      "job.bat",
		Param:     "A",
		Env:       "TESTENV1=!!!WARNING!!!",
		Workspace: "C:\\Go",
		WarnRC:    0,
		WarnStr:   "",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("RCは0のはず.")
	}
	if len(res.Detail) != 0 {
		t.Error("異常終了メッセージが存在しない.")
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "!!!WARNING!!! A" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.NORMAL {
		t.Error("statが正常終了ではない")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_JSジョブが正常に実行できる(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       109,
		JID:       "serviceTask_005",
		Path:      "job.js",
		Param:     "A B",
		Env:       "",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 0 {
		t.Error("ジョブが正常終了するはずなのに異常終了した.")
	}
	if len(res.Detail) > 0 {
		t.Error("ジョブが正常終了するはずなのに、エラーメッセージがある.", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "Argument2=B" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間が無い.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間が無い.")
	}
	if strings.Index(res.Et, "20140330110120") != -1 {
		t.Error("テストのテスト")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログがない。", err.Error())
	}
	defer file.Close()
}

func TestDoJobRequest_標準エラー出力で警告終了するVBSジョブ_RC確認なし(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       110,
		JID:       "serviceTask_006",
		Path:      "stderr.vbs",
		Param:     "!!!WARN",
		Env:       "",
		Workspace: "C:\\Go",
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
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	} else if res.Detail != detailWarnPtn {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "Argument1=!!!WARN" {
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで異常終了するジョブ_閾値と同じ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       111,
		JID:       "serviceTask_008",
		Path:      "job.bat",
		Param:     "X 12",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 12 {
		t.Error("RCは12のはず.")
	}
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	} else if res.Detail != detailErrRC {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数が格納されていない.")
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで異常終了するがチェックしないジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       112,
		JID:       "serviceTask_008",
		Path:      "job.bat",
		Param:     "A 12",
		Env:       "TESTENV1=ENVENVENV",
		Workspace: "C:\\Go",
		WarnRC:    0,
		WarnStr:   "WARN",
		ErrRC:     0,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 12 {
		t.Error("RCは12のはず.")
	}
	if len(res.Detail) != 0 {
		t.Error("異常終了メッセージが存在しない.")
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "ENVENVENV A" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.NORMAL {
		t.Error("statが正常終了ではない")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで異常終了して標準出力で警告終了するジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       113,
		JID:       "serviceTask_008",
		Path:      "job.bat",
		Param:     "X 12",
		Env:       "TESTENV1=WARNING",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 12 {
		t.Error("RCは12のはず.")
	}
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在しない.")
	} else if res.Detail != detailErrRC {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数が格納されていない.")
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_RCで警告終了して標準出力で異常終了するジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       114,
		JID:       "serviceTask_008",
		Path:      "job.bat",
		Param:     "X 11",
		Env:       "TESTENV1=!!!ERROR!!!",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERROR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 11 {
		t.Error("RCは11のはず.")
	}
	if len(res.Detail) == 0 {
		t.Error("異常終了メッセージが存在する.")
	} else if res.Detail != detailErrPtn {
		t.Errorf("想定外のメッセージ - %v", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数が格納されていない.")
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
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_日本語ジョブ(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       115,
		JID:       "serviceTask_009",
		Path:      "あ.bat",
		Param:     "OOO 100",
		Env:       "TESTENV1=!!!ERROR!!!",
		Workspace: "C:\\Go",
		WarnRC:    0,
		WarnStr:   "",
		ErrRC:     0,
		ErrStr:    "",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.RC != 100 {
		t.Error("RCは100のはず.なのに", res.RC)
	}
	if len(res.Detail) != 0 {
		t.Error("異常終了メッセージが存在する.")
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "!!!ERROR!!! OOO" {
		t.Errorf("変数内容が不正.[%s]", res.Var)
	}
	if len(res.St) == 0 {
		t.Error("ジョブ開始時間がない.")
	}
	if len(res.Et) == 0 {
		t.Error("ジョブ終了時間がない.")
	}
	if res.Stat != db.NORMAL {
		t.Error("statが異常終了")
	}
	file, err := os.Open(createJoblogFileName(req, stToLocalTimestamp(res.St), res.NID, res.JID))
	if err != nil {
		t.Error("ジョブログが存在しない.")
	}
	defer file.Close()
}

func TestDoJobRequest_powershellジョブを実行(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       116,
		JID:       "serviceTask_001",
		Path:      "job.ps1",
		Param:     "-a あいうえお -b 123 -z",
		Env:       "TESTENV1=ENVENV",
		Workspace: "C:\\Go",
		WarnRC:    4,
		WarnStr:   "WARN",
		ErrRC:     12,
		ErrStr:    "ERR",
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if req.NID != res.NID {
		t.Error("NIDがリクエストとレスポンスで異なる.")
	}
	if req.JID != res.JID {
		t.Error("IDがリクエストとレスポンスで異なる.")
	}
	if res.RC != 0 {
		t.Error("ジョブが正常終了するはずなのに異常終了した.")
	}
	if len(res.Detail) > 0 {
		t.Error("ジョブが正常終了するはずなのに、エラーメッセージがある.", res.Detail)
	}
	if len(res.Var) == 0 {
		t.Error("変数なし.")
	} else if res.Var != "$z is True" {
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

func TestDoJobRequest_タイムアウト時間を超えたら異常終了する(t *testing.T) {
	req := &message.Request{
		Type:      "request",
		NID:       117,
		JID:       "serviceTask_001",
		Path:      "twosec.bat",
		Param:     "",
		Env:       "",
		Workspace: "C:\\Go",
		WarnRC:    0,
		WarnStr:   "",
		ErrRC:     0,
		ErrStr:    "",
		Timeout:   1,
	}

	stCh := make(chan string, 1)
	res := DoJobRequest(req, conf, stCh)
	close(stCh)
	if res.Stat != db.ABNORMAL {
		t.Error("異常終了するはずが、正常終了している。")
	}
	if len(res.Detail) == 0 {
		t.Error("異常メッセージがセットされていない。")
	}
}
