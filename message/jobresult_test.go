package message

import "testing"

func TestJobResult_実行結果メッセージをパースできる(t *testing.T) {
	message := `{
    "type":"jobresult",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var j JobResult
	err := j.ParseJSON(message)
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	if j.Version != `1.2.3` {
		t.Errorf("取得したversionの値が違います： %s", j.Version)
	}
	if j.NID != 1234 {
		t.Errorf("取得したnidの値が違います： %d", j.NID)
	}
	if j.JID != `job1` {
		t.Errorf("取得したjidの値が違います： %s", j.JID)
	}
	if j.RC != 10 {
		t.Errorf("取得したrcの値が違います： %d", j.RC)
	}
	if j.Stat != 2 {
		t.Errorf("取得したstatの値が違います： %d", j.Stat)
	}
	if j.Var != `somevalue` {
		t.Errorf("取得したvarの値が違います： %s", j.Var)
	}
	if j.St != `2015-03-26 13:21:15.000` {
		t.Errorf("取得したstの値が違います： %s", j.St)
	}
	if j.Et != `2015-03-26 19:21:36.000` {
		t.Errorf("取得したetの値が違います： %s", j.Et)
	}
}

func TestJobResult_JSON文字列としてパースできない場合はエラーが発生する(t *testing.T) {
	message := `
    "type":"jobresult",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var j JobResult
	err := j.ParseJSON(message)

	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestJobResult_typeが違う場合はエラーが発生する(t *testing.T) {
	message := `{
	"type":"response",
	"version":"1.2.3",
	"nid":1234,
	"jid":"job1",
    "rc":10,
	"stat":2,
	"var":"somevalue",
	"st":"2015-03-26 13:21:15.000",
	"et":"2015-03-26 19:21:36.000"
}`

	var j JobResult
	err := j.ParseJSON(message)

	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestJobResult_プロパティ値からJSONメッセージを生成できる(t *testing.T) {
	ServantVersion = "1.2.3"

	var j JobResult
	j.NID = 1234
	j.JID = `job1`
	j.RC = 10
	j.Stat = 2
	j.Var = `somevalue`
	j.St = `2015-03-26 13:21:15.000`
	j.Et = `2015-03-26 19:21:16.000`

	msg, err := j.GenerateJSON()

	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	expect := `{"type":"jobresult","version":"1.2.3","nid":1234,"jid":"job1","rc":10,"stat":2,"var":"somevalue","st":"2015-03-26 13:21:15.000","et":"2015-03-26 19:21:16.000"}`
	if msg != expect {
		t.Error("生成されたJSONメッセージが想定値と違います")
		t.Logf("生成値: %s", msg)
		t.Logf("想定値: %s", expect)
	}
}
