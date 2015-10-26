// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package remote

import (
	"errors"
	"strings"
	"testing"

	"cuto/db"
	"cuto/message"
	"cuto/servant/config"
	"cuto/testutil"
)

func readTestConfig() *config.ServantConfig {
	config.Servant = config.DefaultServantConfig()
	return config.Servant
}

func doTestRequest(req *message.Request, conf *config.ServantConfig, stCh chan<- string) *message.Response {
	res := new(message.Response)
	res.NID = req.NID
	res.RC = 1
	res.Detail = "detail"
	res.JID = req.JID
	res.Stat = db.NORMAL
	res.St = "20150331131524.123456789"
	res.Et = "20150331131525.123456789"

	return res
}

func TestDo_ジョブを実行し結果を送信できる(t *testing.T) {
	reqMsg := `{
	"type":"request",
	"varsion":"1.2.3",
	"nid":1234,
	"jid":"001",
	"path":"C:\\work\\test.bat",
	"param":"test",
	"env":"testenv=val",
	"workspace": "C:\\work",
	"warnrc":4,
	"warnstr":"warn",
	"errrc":12,
	"errstr":"error"
}`

	conf := readTestConfig()
	message.ServantVersion = "2.3.4"

	conn := testutil.NewConnStub()
	session := Session{Conn: conn, Body: reqMsg, doJobRequest: doTestRequest}
	session.startHeartbeat()
	err := session.Do(conf)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	expected := `{"type":"response","version":"2.3.4","nid":1234,"jid":"001","rc":1,"stat":1,"detail":"detail","var":"","st":"20150331131524.123456789","et":"20150331131525.123456789"}`
	expected += "\n"
	if conn.WriteStr != expected {
		t.Errorf("送信されたジョブ実行結果が間違っています。")
		t.Logf("想定値: %s", expected)
		t.Logf("実績値: %s", conn.WriteStr)
	}
}

func TestDo_パースできないリクエストメッセージが来たらエラー(t *testing.T) {
	reqMsg := `notjson`

	conf := readTestConfig()

	conn := testutil.NewConnStub()
	session := Session{Conn: conn, Body: reqMsg, doJobRequest: doTestRequest}
	session.startHeartbeat()
	err := session.Do(conf)
	if err == nil {
		t.Error("エラーが発生していない。")
	}

	if conn.WriteStr != "" {
		t.Errorf("想定外のメッセージが書き込まれた: %s", conn.WriteStr)
	}
}

func TestDo_使用不可能な変数が使用されたら異常終了のレスポンスを返す(t *testing.T) {
	reqMsg := `{
	"type":"request",
	"varsion":"1.2.3",
	"nid":1234,
	"jid":"001",
	"path":"C:\\work\\test.bat",
	"param":"$SSJOBNET:ID$",
	"env":"testenv=val",
	"workspace": "C:\\work",
	"warnrc":4,
	"warnstr":"warn",
	"errrc":12,
	"errstr":"error"
}`

	conf := readTestConfig()

	conn := testutil.NewConnStub()
	session := Session{Conn: conn, Body: reqMsg, doJobRequest: doTestRequest}
	session.startHeartbeat()
	err := session.Do(conf)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	if !strings.Contains(conn.WriteStr, `"stat":9`) || !strings.Contains(conn.WriteStr, "Undefined variable") {
		t.Errorf("想定外のメッセージが書き込まれた: %s", conn.WriteStr)
	}
}

func TestDo_コネクションオブジェクトへのWriteに失敗したらエラー(t *testing.T) {
	reqMsg := `{
	"type":"request",
	"varsion":"1.2.3",
	"nid":1234,
	"jid":"001",
	"path":"C:\\work\\test.bat",
	"param":"test",
	"env":"testenv=val",
	"workspace": "C:\\work",
	"warnrc":4,
	"warnstr":"warn",
	"errrc":12,
	"errstr":"error"
}`

	conf := readTestConfig()

	conn := testutil.NewConnStub()
	conn.WriteErr = errors.New("testerror")
	session := Session{Conn: conn, Body: reqMsg, doJobRequest: doTestRequest}
	session.startHeartbeat()
	err := session.Do(conf)
	if err == nil {
		t.Error("エラーが発生していない。")
	}

	if conn.WriteStr != "" {
		t.Errorf("想定外のメッセージが書き込まれた: %s", conn.WriteStr)
	}
}
