package message

import "testing"

func TestResponse_実行結果メッセージをパースできる(t *testing.T) {
	message := `{
    "type":"response",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"detail":"something",
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var res Response
	err := res.ParseJSON(message)
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	if res.Version != `1.2.3` {
		t.Errorf("取得したversionの値が違います： %s", res.Version)
	}
	if res.NID != 1234 {
		t.Errorf("取得したnidの値が違います： %d", res.NID)
	}
	if res.JID != `job1` {
		t.Errorf("取得したjidの値が違います： %s", res.JID)
	}
	if res.RC != 10 {
		t.Errorf("取得したrcの値が違います： %d", res.RC)
	}
	if res.Stat != 2 {
		t.Errorf("取得したstatの値が違います： %d", res.Stat)
	}
	if res.Detail != `something` {
		t.Errorf("取得したdetailの値が違います： %s", res.Detail)
	}
	if res.Var != `somevalue` {
		t.Errorf("取得したvarの値が違います： %s", res.Var)
	}
	if res.St != `2015-03-26 13:21:15.000` {
		t.Errorf("取得したstの値が違います： %s", res.St)
	}
	if res.Et != `2015-03-26 19:21:36.000` {
		t.Errorf("取得したetの値が違います： %s", res.Et)
	}
}

func TestResponse_JSON文字列としてパースできない場合はエラーとする(t *testing.T) {
	message := `
    "type":"response",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"detail":"something",
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var res Response
	err := res.ParseJSON(message)

	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestResponse_メッセージタイプが違う場合はエラーとする(t *testing.T) {
	message := `{
	"type":"somthingelse",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"detail":"something",
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var res Response
	err := res.ParseJSON(message)

	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestResponse_プロパティ値からJSONメッセージを生成できる(t *testing.T) {
	ServantVersion = "1.2.3"

	var res Response
	res.NID = 1234
	res.JID = `job1`
	res.RC = 10
	res.Stat = 2
	res.Detail = `something`
	res.Var = `somevalue`
	res.St = `2015-03-26 13:21:15.000`
	res.Et = `2015-03-26 19:21:16.000`

	msg, err := res.GenerateJSON()

	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	expect := `{"type":"response","version":"1.2.3","nid":1234,"jid":"job1","rc":10,"stat":2,"detail":"something","var":"somevalue","st":"2015-03-26 13:21:15.000","et":"2015-03-26 19:21:16.000"}`
	if msg != expect {
		t.Error("生成されたJSONメッセージが想定値と違います")
		t.Logf("生成値: %s", msg)
		t.Logf("想定値: %s", expect)
	}
}
