// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package remote

import (
	"testing"

	"cuto/db"

	"cuto/message"
	"cuto/servant/config"
)

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

	conf := config.ReadConfig()

	conn := new(testConn)
	session := Session{Conn: conn, Body: reqMsg, doRequest: doTestRequest}
	session.Do(conf)

	expected := `{"type":"response","nid":1234,"jid":"001","rc":1,"stat":1,"detail":"detail","var":"","st":"20150331131524.123456789","et":"20150331131525.123456789"}`
	if conn.Written != expected {
		t.Errorf("送信されたジョブ実行結果が間違っています。")
		t.Logf("想定値: %s", expected)
		t.Logf("実績値: %s", conn.Written)
	}
}
