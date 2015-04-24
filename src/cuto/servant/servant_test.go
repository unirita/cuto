package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"cuto/servant/config"
	"cuto/servant/remote"
	"cuto/testutil"
)

func waitEventLoopEnd(sigCh <-chan os.Signal, sq <-chan *remote.Session, endCh chan<- struct{}) {
	eventLoop(sigCh, sq)
	endCh <- struct{}{}
}

func TestRun_メッセージ受信開始に失敗したらエラー(t *testing.T) {
	config.Servant = config.DefaultServantConfig()
	config.Servant.Sys.BindPort = 65536
	_, err := Run()
	if err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestEventLoop_SIGINTシグナルを受信するとループが終了する(t *testing.T) {
	isTest = true

	sigCh := make(chan os.Signal)
	sq := make(chan *remote.Session)
	endCh := make(chan struct{})
	go waitEventLoopEnd(sigCh, sq, endCh)

	sigCh <- syscall.SIGINT
	select {
	case <-endCh:
		// 問題なし
	case <-time.After(time.Millisecond * 100):
		t.Errorf("ループが終了しない。")
	}
}

func TestEventLoop_ハングアップしてもループが終了しない(t *testing.T) {
	isTest = true

	sigCh := make(chan os.Signal)
	sq := make(chan *remote.Session)
	endCh := make(chan struct{})
	go waitEventLoopEnd(sigCh, sq, endCh)

	sigCh <- syscall.SIGHUP
	select {
	case <-endCh:
		t.Errorf("ループが終了した。")
	case <-time.After(time.Millisecond * 100):
		sigCh <- syscall.SIGINT
	}
}

func TestEventLoop_セッションキューに挿入したセッションが実行される(t *testing.T) {
	isTest = true

	// あえてエラーが発生するSession.Bodyをセットし、doJobRequestまでロジックを進ませない
	body := `{
	"type":"request",
	"nid":1234,
	"jid":"001",
	"path":"test.bat",
	"param":"$SSJOBNET:ID$",
	"env":"testenv=val",
	"workspace": "",
	"warnrc":4,
	"warnstr":"warn",
	"errrc":12,
	"errstr":"error"
}`

	conn := testutil.NewConnStub()
	session := remote.NewSession(conn, body)

	sigCh := make(chan os.Signal)
	sq := make(chan *remote.Session)
	go eventLoop(sigCh, sq)

	sq <- session

	// Session実行ごルーチンの終了待ちのためにwait
	time.Sleep(100 * time.Millisecond)
	if conn.WriteStr == "" {
		t.Error("セッションが実行されていない。")
	}
}
