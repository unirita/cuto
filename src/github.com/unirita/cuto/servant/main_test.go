package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/unirita/cuto/servant/config"
	"github.com/unirita/cuto/testutil"
)

func getTestDataDir() string {
	return "_testdata"
}

func TestRealMain_バージョン確認ができる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	args := new(arguments)
	args.v = true

	c.Start()
	realMain(args)
	out := c.Stop()

	if !strings.Contains(out, Version) {
		t.Error("バージョンが出力されていない。")
	}
}

func TestRealMain_設定ファイルから設定がロードされた上で内容にエラーがあればリターンコードrc_errorを返す(t *testing.T) {
	const s = os.PathSeparator
	var configFile string
	if runtime.GOOS == "windows" {
		configFile = "error.ini"
	} else {
		configFile = "error_l.ini"
	}
	config.FilePath = filepath.Join(getTestDataDir(), configFile)

	args := new(arguments)
	rc := realMain(args)

	if config.Servant.Sys.BindPort != 65536 {
		t.Error("取得した設定値が想定と違っている。")
	}
	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestRealMain_ロガー初期化でのエラー発生時にリターンコードrc_errorを返す(t *testing.T) {
	const s = os.PathSeparator
	var configFile string
	if runtime.GOOS == "windows" {
		configFile = "logerror.ini"
	} else {
		configFile = "logerror_l.ini"
	}
	config.FilePath = filepath.Join(getTestDataDir(), configFile)

	args := new(arguments)
	rc := realMain(args)

	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestRealMain_Run関数でのエラー発生時にリターンコードrc_errorを返す(t *testing.T) {
	var configFile string
	const s = os.PathSeparator
	if runtime.GOOS == "windows" {
		configFile = "binderror.ini"
	} else {
		configFile = "binderror_l.ini"
	}
	config.FilePath = filepath.Join(getTestDataDir(), configFile)

	args := new(arguments)
	rc := realMain(args)

	if rc != rc_error {
		t.Errorf("リターンコード[%d]が想定値と違っている。", rc)
	}
}

func TestFetchArgs_実行時引数を取得できる(t *testing.T) {
	os.Args = append(os.Args, "-v")
	args := fetchArgs()
	if !args.v {
		t.Error("バージョン出力オプションが取得できていない。")
	}
}
