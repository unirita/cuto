package message

import (
	"testing"
)

func TestJobCheck_ジョブ正常終了確認メッセージをパースできる(t *testing.T) {
	message := `{
    "type":"jobcheck",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1"
}`

	var j JobCheck
	err := j.ParseJSON(message)
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}
	if j.Version != "1.2.3" {
		t.Errorf("取得したversionの値が違います： %s", j.Version)
	}
	if j.NID != 1234 {
		t.Errorf("取得したnidの値が違います： %d", j.NID)
	}
	if j.JID != `job1` {
		t.Errorf("取得したjidの値が違います： %s", j.JID)
	}
}

func TestJobCheck_JSONとしてパースできない文字列の場合はエラーが発生する(t *testing.T) {
	message := `
    "type":"jobcheck",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1"
}`

	var j JobCheck
	err := j.ParseJSON(message)
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestJobCheck_typeが間違っている場合はエラーが発生する(t *testing.T) {
	message := `{
    "type":"request",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1"
}`

	var j JobCheck
	err := j.ParseJSON(message)
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestJobCheck_オブジェクトからJSONメッセージを生成できる(t *testing.T) {
	MasterVersion = "1.2.3"

	var j JobCheck
	j.NID = 1234
	j.JID = `job1`

	msg, err := j.GenerateJSON()
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	expect := `{"type":"jobcheck","version":"1.2.3","nid":1234,"jid":"job1"}`
	if msg != expect {
		t.Error("生成されたJSONメッセージが想定値と違います")
		t.Logf("生成値: %s", msg)
		t.Logf("想定値: %s", expect)
	}
}
