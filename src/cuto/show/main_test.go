// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"os/exec"

	"testing"

	"cuto/console"
	"cuto/db"
	"cuto/show/gen"
	"cuto/testutil"
)

// テスト結果確認用の一時出力ファイル。
const output_file string = "showtest.log"

var (
	confPath = setConfigPath()
	confFile = setConfigFile()
)
var diffCommand = getDiff()

func getDiff() string {
	if runtime.GOOS == "windows" {
		return "fc"
	} else {
		return "diff"
	}
}

func setConfigPath() string {
	return fmt.Sprintf("%v%c%v%c%v%c%v", os.Getenv("GOPATH"),
		os.PathSeparator, "test", os.PathSeparator, "cuto", os.PathSeparator, "show")
}
func setConfigFile() string {
	return fmt.Sprintf("%v%c%v", confPath, os.PathSeparator, "show_test.ini")
}

func init() {
	os.Chdir(confPath) // 設定ファイル内を固定するため、作業フォルダを固定する。
}

func vefiry_stdout(output_file, vefiry_file string) error {
	params := []string{output_file, vefiry_file}
	cmd := exec.Command(diffCommand, params...)
	err := cmd.Run()
	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("不正な戻り値[%v]が返りました。", s.ExitStatus())
			} else {
				return fmt.Errorf("不正な結果です。")
			}
		} else {
			return fmt.Errorf("不正な結果です。")
		}
	}
	return nil
}

func TestRealMain_1日分のジョブネットを表示(t *testing.T) {
	vefiry_file := "showtest_verify1.txt"
	arg := &arguments{
		jobnet: "",
		from:   "20150415",
		to:     "20150415",
		status: "",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_OK {
		t.Errorf("戻り値[%v]が返るはずが、[%v]が返りました。", rc_OK, ret)
	}
	cout := co.Stop()
	cerr := ce.Stop()
	if len(cerr) > 0 {
		t.Errorf("エラーが出力されています。 - %v", cerr)
	}
	if _, exist := os.Stat(output_file); exist != nil {
		os.Remove(output_file)
	}
	file, _ := os.OpenFile(output_file, os.O_CREATE|os.O_WRONLY, 0666)
	file.WriteString(cout)
	file.Close()

	err := vefiry_stdout(output_file, vefiry_file)
	if err != nil {
		t.Errorf("不正な出力結果です。 - %v", err)
	}
}

func TestRealMain_0件のジョブネットを表示(t *testing.T) {
	arg := &arguments{
		jobnet: "JNET",
		from:   "",
		to:     "",
		status: "normal",
		format: "csv",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_NOTHING {
		t.Errorf("戻り値[%v]が返るはずが、[%v]が返りました。", rc_NOTHING, ret)
	}
	cout := co.Stop()
	cerr := ce.Stop()
	if len(cerr) > 0 {
		t.Errorf("エラーが出力されています。 - %v", cerr)
	}
	if len(cout) > 0 {
		t.Errorf("出力されないはずが、何か出力されました。 - %v", cout)
	}
}

func TestRealMain_ヘルプを表示(t *testing.T) {
	arg := &arguments{
		help:   true,
		jobnet: "",
		from:   "20150416",
		to:     "20150416",
		status: "",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_OK {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_OK, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力が出力されています。 - %v", cout)
	}
	if cerr != console.USAGE_SHOW {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", cerr)
	}
}

func TestRealMain_バージョン情報を表示(t *testing.T) {
	arg := &arguments{
		v:      true,
		jobnet: "",
		from:   "20150416",
		to:     "20150416",
		status: "",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_OK {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_OK, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力が出力されています。 - %v", cout)
	}
	if !strings.Contains(cerr, Version) {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", cerr)
	}
}

func TestRealMain_不正なFROM(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150416X",
		to:     "20150416",
		status: "",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_PARMERR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_PARMERR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, console.USAGE_SHOW) {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_不正なTO(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150416",
		to:     "20150416X",
		status: "",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_PARMERR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_PARMERR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, console.USAGE_SHOW) {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_不正なStatus(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150416",
		to:     "20150416",
		status: "abc",
		format: "",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_PARMERR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_PARMERR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, console.USAGE_SHOW) {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_不正なFormat(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150416",
		to:     "20150416",
		status: "",
		format: "abc",
		config: confFile,
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_PARMERR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_PARMERR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, console.USAGE_SHOW) {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_不正な設定ファイル(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150416",
		to:     "20150416",
		status: "",
		format: "",
		config: "error.ini",
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_PARMERR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_PARMERR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, "error.ini") {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_RuntimeErrorX(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150414",
		to:     "20150416",
		status: "",
		format: "",
		config: fmt.Sprintf("%v%c%v", confPath, os.PathSeparator, "show_testX.ini"),
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_ERROR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_ERROR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, "AN INTERNAL ERROR OCCURRED.") {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestRealMain_RuntimeErrorY(t *testing.T) {
	arg := &arguments{
		jobnet: "",
		from:   "20150414",
		to:     "20150416",
		status: "",
		format: "",
		config: fmt.Sprintf("%v%c%v", confPath, os.PathSeparator, "show_testY.ini"),
	}
	ce := testutil.NewStderrCapturer()
	ce.Start()
	co := testutil.NewStdoutCapturer()
	co.Start()

	ret := realMain(arg)
	if ret != rc_ERROR {
		t.Errorf("戻り値[%v]が返るべきところ、[%v]が返りました。", rc_ERROR, ret)
	}

	cout := co.Stop()
	cerr := ce.Stop()
	if len(cout) > 0 {
		t.Errorf("不正な標準出力です。 - %v", cout)
	}
	if !strings.Contains(cerr, "AN INTERNAL ERROR OCCURRED.") {
		t.Errorf("不正な標準エラー出力です。 - %v", cerr)
	}
}

func TestGetStatusType_指定毎に返るステータスを確認(t *testing.T) {
	s, err := getStatusType("running")
	if err != nil {
		t.Errorf("成功する予定が、エラーになった。 - %v", err)
	} else if s != db.RUNNING {
		t.Errorf("ステータス[%v]が返る予定が、[%v]が返った。", db.RUNNING, s)
	}

	s, err = getStatusType("normal")
	if err != nil {
		t.Errorf("成功する予定が、エラーになった。 - %v", err)
	} else if s != db.NORMAL {
		t.Errorf("ステータス[%v]が返る予定が、[%v]が返った。", db.NORMAL, s)
	}

	s, err = getStatusType("abnormal")
	if err != nil {
		t.Errorf("成功する予定が、エラーになった。 - %v", err)
	} else if s != db.ABNORMAL {
		t.Errorf("ステータス[%v]が返る予定が、[%v]が返った。", db.ABNORMAL, s)
	}

	s, err = getStatusType("warn")
	if err != nil {
		t.Errorf("成功する予定が、エラーになった。 - %v", err)
	} else if s != db.WARN {
		t.Errorf("ステータス[%v]が返る予定が、[%v]が返った。", db.WARN, s)
	}

	s, err = getStatusType("")
	if err != nil {
		t.Errorf("成功する予定が、エラーになった。 - %v", err)
	} else if s != -1 {
		t.Errorf("ステータス[%v]が返る予定が、[%v]が返った。", -1, s)
	}

	s, err = getStatusType("X")
	if err == nil {
		t.Error("失敗する予定が、エラーが返らなかった。")
	}
}

func TestGetSeparatorType_指定毎に返るジェネレーターを確認(t *testing.T) {
	g := getSeparatorType("")
	switch g.(type) {
	case gen.JsonGenerator:
	default:
		t.Error("JsonGeneratorになるべきところ、異なる型が返った。")
	}

	g = getSeparatorType("json")
	switch g.(type) {
	case gen.JsonGenerator:
	default:
		t.Error("JsonGeneratorになるべきところ、異なる型が返った。")
	}

	g = getSeparatorType("csv")
	switch g.(type) {
	case gen.CsvGenerator:
	default:
		t.Error("CsvGeneratorになるべきところ、異なる型が返った。")
	}

	g = getSeparatorType("X")
	if g != nil {
		t.Error("誤った指定をしたにもかかわらず、nilが返らない。")
	}
}

func TestShowVersion_バージョン表示(t *testing.T) {
	c := testutil.NewStderrCapturer()
	c.Start()

	showVersion()

	output := c.Stop()
	if output != fmt.Sprintf("%v show-utility version.\n", Version) {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestShowUsage_Usage表示(t *testing.T) {
	c := testutil.NewStderrCapturer()
	c.Start()

	showUsage()

	output := c.Stop()
	if output != console.USAGE_SHOW {
		t.Errorf("stderrへの出力値[%s]が想定と違います。", output)
	}
}

func TestFetchArgs_実行時引数のフェッチ(t *testing.T) {
	fetchArgs()
}
