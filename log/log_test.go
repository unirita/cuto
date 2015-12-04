package log

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cihub/seelog"

	"github.com/unirita/cuto/testutil"
	"github.com/unirita/cuto/util"
)

var lockName = getLockName()

func getLockName() string {
	if runtime.GOOS == "windows" {
		return "log_test"
	} else {
		return "test.lock"
	}
}

func initForTest() {
	config := `
<seelog type="sync" minlevel="trace">
    <outputs formatid="common">
        <console />
    </outputs>
    <formats>
        <format id="common" format="2015-04-01 12:34:56.789 [%LEV] %Msg%n"/>
    </formats>
</seelog>`
	logger, _ := seelog.LoggerFromConfigAsString(config)
	locker, _ = util.InitLock(lockName)
	seelog.ReplaceLogger(logger)
	valid = true
}

func TestInit_ログレベルの指定が不正な場合にエラーが発生する(t *testing.T) {
	err := Init("", "test", "", "invalid", 1000, 2, 1)
	if err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestInit_初期化処理終了時点でログファイルがあらかじめ生成される(t *testing.T) {
	testDir := filepath.Join(testutil.GetBaseDir(), "log", "_tmp")
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0777)

	Init(testDir, "test", "", "info", 1000, 2, 1)
	expectedFile := filepath.Join(testDir, "test.log")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Error("Log file must be created, but it was not.")
		t.Logf("filepath: %s", expectedFile)
	}
}

func TestTrace_traceレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Trace("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [TRC] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestTrace_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Trace("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestDebug_debugレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Debug("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [DBG] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestDebug_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Debug("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestInfo_infoレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Info("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [INF] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestInfo_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Info("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestWarn_warnレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Warn("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [WRN] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestWarn_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Warn("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestError_errorレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Error("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [ERR] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestError_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Error("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestCritical_criticalレベルのメッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	initForTest()
	defer Term()

	c.Start()
	Critical("testmessage")
	output := c.Stop()

	if output != "2015-04-01 12:34:56.789 [CRT] testmessage\n" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}

func TestCritical_Initされていない時は何も出力しない(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	valid = false

	c.Start()
	Critical("testmessage")
	output := c.Stop()

	if output != "" {
		t.Errorf("出力されたメッセージ[%s]が想定と異なる。", output)
	}
}
