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

	"path/filepath"

	"cuto/db"
	"cuto/message"
	"cuto/servant/config"
)

var conf *config.ServantConfig
var testJobPath string

// ジョブログなどの掃除
func init() {
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("Location Asia/Tokyo is not found.")
	}

	testJobPath = filepath.Join(os.Getenv("GOPATH"), "test", "cuto", "servant", "job")
	err = os.Chdir(testJobPath)
	config.RootPath = testJobPath
	if err != nil {
		panic(err.Error())
	}
	configPath := filepath.Join(testJobPath, "servant_l.ini")
	conf = config.ReadConfig(configPath)
	os.RemoveAll(conf.Dir.JoblogDir)
	err = os.Mkdir(conf.Dir.JoblogDir, 0755)
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
	joblogDir := filepath.Join(conf.Dir.JoblogDir, st[:8])
	return fmt.Sprintf("%v%c%v.%v.%v.%v.log", joblogDir, os.PathSeparator, nID, job, jID, st)
}

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
	req.Path = filepath.Join(testJobPath, "jobscript", "job.csh")

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
