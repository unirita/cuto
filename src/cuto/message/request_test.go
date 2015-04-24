package message

import (
	"os"
	"testing"
)

func TestRequest_実行命令メッセージをパースできる(t *testing.T) {
	message := `{
    "type":"request",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1",
    "path":"C:\\work\\test.bat",
    "param":"test",
    "env":"testenv",
    "workspace":"C:\\work",
    "warnrc":10,
    "warnstr":"warn",
    "errrc":20,
    "errstr":"err",
    "timeout":60
}`

	var req Request
	err := req.ParseJSON(message)
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	if req.Version != "1.2.3" {
		t.Errorf("取得したversionの値が違います： %s", req.Version)
	}
	if req.NID != 1234 {
		t.Errorf("取得したnidの値が違います： %d", req.NID)
	}
	if req.JID != `job1` {
		t.Errorf("取得したjidの値が違います： %s", req.JID)
	}
	if req.Path != `C:\work\test.bat` {
		t.Errorf("取得したpathの値が違います： %s", req.Path)
	}
	if req.Param != `test` {
		t.Errorf("取得したparamの値が違います： %s", req.Param)
	}
	if req.Env != `testenv` {
		t.Errorf("取得したenvの値が違います： %s", req.Env)
	}
	if req.Workspace != `C:\work` {
		t.Errorf("取得したWorkspaceの値が違います： %s", req.Workspace)
	}
	if req.WarnRC != 10 {
		t.Errorf("取得したwarnrcの値が違います： %d", req.WarnRC)
	}
	if req.WarnStr != `warn` {
		t.Errorf("取得したwarnstrの値が違います： %s", req.WarnStr)
	}
	if req.ErrRC != 20 {
		t.Errorf("取得したerrrcの値が違います： %d", req.ErrRC)
	}
	if req.ErrStr != `err` {
		t.Errorf("取得したerrstrの値が違います： %s", req.ErrStr)
	}
	if req.Timeout != 60 {
		t.Errorf("取得したtimeoutの値が違います： %d", req.Timeout)
	}
}

func TestRequest_JSON文字列としてパースできない場合はエラーとする(t *testing.T) {
	message := `
    "type":"request",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1",
    "path":"C:\\work\\test.bat",
    "param":"test",
    "env":"testenv",
    "workspace":"C:\\work",
    "warnrc":10,
    "warnstr":"warn",
    "errrc":20,
    "errstr":"err",
    "timeout":60
}`

	var req Request
	err := req.ParseJSON(message)
	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestRequest_メッセージタイプが違う場合はエラーとする(t *testing.T) {
	message := `{
	"type":"somthingelse",
    "version":"1.2.3",
    "nid":1234,
    "jid":"job1",
    "path":"C:\\work\\test.bat",
    "param":"test",
    "env":"testenv",
    "workspace":"C:\\work",
    "warnrc":10,
    "warnstr":"warn",
    "errrc":20,
    "errstr":"err",
    "timeout":60
}`

	var req Request
	err := req.ParseJSON(message)
	if err == nil {
		t.Error("発生すべきエラーが発生しませんでした。")
	}
}

func TestRequest_プロパティ値からJSONメッセージを生成できる(t *testing.T) {
	MasterVersion = "1.2.3"

	var req Request
	req.NID = 1234
	req.JID = `job1`
	req.Path = `C:\work\test.bat`
	req.Param = `test`
	req.Env = `testenv`
	req.Workspace = `C:\work`
	req.WarnRC = 10
	req.WarnStr = `warn`
	req.ErrRC = 20
	req.ErrStr = `err`
	req.Timeout = 60

	msg, err := req.GenerateJSON()
	if err != nil {
		t.Fatalf("想定外のエラーが発生しました: %s", err)
	}

	expect := `{"type":"request","version":"1.2.3","nid":1234,"jid":"job1","path":"C:\\work\\test.bat","param":"test","env":"testenv","workspace":"C:\\work","warnrc":10,"warnstr":"warn","errrc":20,"errstr":"err","timeout":60}`
	if msg != expect {
		t.Error("生成されたJSONメッセージが想定値と違います")
		t.Logf("生成値: %s", msg)
		t.Logf("想定値: %s", expect)
	}
}

func TestExpandMasterVars_メッセージ内の変数を展開出来る(t *testing.T) {
	AddSysValue("JOBNET", "ID", "123456")
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	AddJobValue("test", res)

	req := new(Request)
	req.Path = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Param = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Env = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Workspace = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`

	err := req.ExpandMasterVars()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	expected := `$MSJOBNET:ID$ testenv $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	if req.Path != expected {
		t.Errorf("変数展開後のPathの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Path)
	}
	if req.Workspace != expected {
		t.Errorf("変数展開後のWorkspaceの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Workspace)
	}
	expected = `123456 testenv 1 $SSROOT$ $SETEST$ $SJtest:RC$`
	if req.Param != expected {
		t.Errorf("変数展開後のParamの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Param)
	}
	if req.Env != expected {
		t.Errorf("変数展開後のEnvの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Env)
	}
}

func TestExpandMasterVars_エラー発生時は元の値を維持する(t *testing.T) {
	AddSysValue("JOBNET", "ID", "123456")
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	AddJobValue("test", res)

	req := new(Request)
	req.Path = `$MSJOBNET:ID$ $METEST$ $MJtest1:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Param = `$MSJOBNET:ID$ $METEST$ $MJtest1:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Env = `$MSJOBNET:ID$ $METEST$ $MJtest1:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Workspace = `$MSJOBNET:ID$ $METEST$ $MJtest1:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`

	err := req.ExpandMasterVars()
	if err == nil {
		t.Fatal("エラーが発生していない。")
	}

	expected := `$MSJOBNET:ID$ $METEST$ $MJtest1:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	if req.Path != expected {
		t.Errorf("変数展開後のPathの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Path)
	}
	if req.Workspace != expected {
		t.Errorf("変数展開後のWorkspaceの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Workspace)
	}
	if req.Param != expected {
		t.Errorf("変数展開後のParamの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Param)
	}
	if req.Env != expected {
		t.Errorf("変数展開後のEnvの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Env)
	}
}

func TestExpandServantVars_メッセージ内の変数を展開出来る(t *testing.T) {
	AddSysValue("ROOT", "", `C:\cute`)
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	AddJobValue("test", res)

	req := new(Request)
	req.Path = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Param = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Env = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`
	req.Workspace = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSROOT$ $SETEST$ $SJtest:RC$`

	err := req.ExpandServantVars()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	expected := `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ C:\cute testenv $SJtest:RC$`
	if req.Path != expected {
		t.Errorf("変数展開後のPathの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Path)
	}
	if req.Workspace != expected {
		t.Errorf("変数展開後のWorkspaceの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Workspace)
	}
	if req.Param != expected {
		t.Errorf("変数展開後のParamの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Param)
	}
	if req.Env != expected {
		t.Errorf("変数展開後のEnvの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Env)
	}
}

func TestExpandServantVars_エラー発生時は元の値を維持する(t *testing.T) {
	AddSysValue("ROOT", "", `C:\cute`)
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	AddJobValue("test", res)

	req := new(Request)
	req.Path = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSNOEXST$ $SETEST$ $SJtest:RC$`
	req.Param = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSNOEXST$ $SETEST$ $SJtest:RC$`
	req.Env = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSNOEXST$ $SETEST$ $SJtest:RC$`
	req.Workspace = `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSNOEXST$ $SETEST$ $SJtest:RC$`

	err := req.ExpandServantVars()
	if err == nil {
		t.Fatal("エラーが発生していない。")
	}

	expected := `$MSJOBNET:ID$ $METEST$ $MJtest:RC$ $SSNOEXST$ $SETEST$ $SJtest:RC$`
	if req.Path != expected {
		t.Errorf("変数展開後のPathの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Path)
	}
	if req.Workspace != expected {
		t.Errorf("変数展開後のWorkspaceの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Workspace)
	}
	if req.Param != expected {
		t.Errorf("変数展開後のParamの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Param)
	}
	if req.Env != expected {
		t.Errorf("変数展開後のEnvの値が想定と違っている。")
		t.Logf("想定値：%s", expected)
		t.Logf("実績値：%s", req.Env)
	}
}
