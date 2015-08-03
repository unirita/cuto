package message

import (
	"os"
	"testing"
)

func TestNewVariable_変数オブジェクトを生成できる(t *testing.T) {
	v := NewVariable(`$METEST$`)
	if v == nil {
		t.Fatal("オブジェクト生成に失敗した。")
	}
	if v.String() != `$METEST$` {
		t.Errorf("文字列表現[%s]が想定と違っている", v.String())
	}
	if v.Place != 'M' {
		t.Errorf("Placeの値[%c]が想定と違っている。", v.Place)
	}
	if v.Kind != 'E' {
		t.Errorf("Kindの値[%c]が想定と違っている。", v.Kind)
	}
	if v.Name != `TEST` {
		t.Errorf("Nameの値[%s]が想定と違っている。", v.Name)
	}
	if v.Tag != `` {
		t.Errorf("Tagの値[%s]が想定と違っている。", v.Tag)
	}
}

func TestNewVariable_変数名が明らかに短い場合nilを返す(t *testing.T) {
	v := NewVariable(`$ME$`)
	if v != nil {
		t.Error("nilを返さなかった。")
	}
}

func TestNewVariable_場所識別子が定義外の場合nilを返す(t *testing.T) {
	v := NewVariable(`$AETEST$`)
	if v != nil {
		t.Error("nilを返さなかった。")
	}
}

func TestNewVariable_種別識別子が定義外の場合nilを返す(t *testing.T) {
	v := NewVariable(`$MATEST$`)
	if v != nil {
		t.Error("nilを返さなかった。")
	}
}

func TestNewVariable_変数オブジェクトを生成できる_タグ付き(t *testing.T) {
	v := NewVariable(`$MJTEST:ID$`)
	if v == nil {
		t.Fatal("オブジェクト生成に失敗した。")
	}
	if v.String() != `$MJTEST:ID$` {
		t.Errorf("文字列表現[%s]が想定と違っている", v.String())
	}
	if v.Place != 'M' {
		t.Errorf("Placeの値[%c]が想定と違っている。", v.Place)
	}
	if v.Kind != 'J' {
		t.Errorf("Kindの値[%c]が想定と違っている。", v.Kind)
	}
	if v.Name != `TEST` {
		t.Errorf("Nameの値[%s]が想定と違っている。", v.Name)
	}
	if v.Tag != `ID` {
		t.Errorf("Tagの値[%s]が想定と違っている。", v.Tag)
	}
}

func TestExpand_システム変数の値を取得できる(t *testing.T) {
	AddSysValue(`JOBNET`, `ID`, `123456`)
	AddSysValue(`ROOT`, ``, `C:\cute`)

	v := NewVariable("$MSJOBNET:ID$")
	val, err := v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != `123456` {
		t.Errorf("取得したJOBNET:IDの値[%s]が想定と違っている。", val)
	}

	v = NewVariable("$SSROOT$")
	val, err = v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != `C:\cute` {
		t.Errorf("取得したSSROOTの値[%s]が想定と違っている。", val)
	}
}

func TestExpand_未定義のシステム変数を参照したらエラー(t *testing.T) {
	v := NewVariable("$MSUNDEF$")
	_, err := v.Expand()
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestExpand_環境変数の値を取得できる(t *testing.T) {
	v := NewVariable("$METEST$")
	os.Setenv("TEST", "testenv")

	val, err := v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "testenv" {
		t.Errorf("取得した値[%s]が想定と違っている。", val)
	}
}

func TestExpand_ジョブネットワーク変数の値を取得できる(t *testing.T) {
	res := new(Response)
	res.JID = "JOB1"
	res.RC = 1
	res.St = "2015-01-01 12:34:56.789"
	res.Et = "2015-01-01 13:00:00.000"
	res.Var = "testout"
	AddJobValue("test", res)

	v := NewVariable("$MJtest:ID$")
	val, err := v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "JOB1" {
		t.Errorf("取得したRCの値[%s]が想定と違っている。", val)
	}

	v = NewVariable("$MJtest:RC$")
	val, err = v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "1" {
		t.Errorf("取得したRCの値[%s]が想定と違っている。", val)
	}

	v = NewVariable("$MJtest:SD$")
	val, err = v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "$ST20150101123456.789$" {
		t.Errorf("取得したSDの値[%s]が想定と違っている。", val)
	}

	v = NewVariable("$MJtest:ED$")
	val, err = v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "$ST20150101130000.000$" {
		t.Errorf("取得したEDの値[%s]が想定と違っている。", val)
	}

	v = NewVariable("$MJtest:OUT$")
	val, err = v.Expand()
	if err != nil {
		t.Fatalf("想定外のエラーが発生[%s]", err)
	}
	if val != "testout" {
		t.Errorf("取得したOUTの値[%s]が想定と違っている。", val)
	}
}

func TestExpand_未定義のジョブネットワーク変数を参照したらエラー(t *testing.T) {
	res := new(Response)
	res.RC = 1
	res.St = "2015-01-01 12:34:56.789"
	res.Et = "2015-01-01 13:00:00.000"
	res.Var = "testout"
	AddJobValue("test", res)

	v := NewVariable("$MJtest2:RC$")
	_, err := v.Expand()
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestExpand_servantではジョブネットワーク変数を参照したらエラー(t *testing.T) {
	res := new(Response)
	res.RC = 1
	res.St = "2015-01-01 12:34:56.789"
	res.Et = "2015-01-01 13:00:00.000"
	res.Var = "testout"
	AddJobValue("test", res)

	v := NewVariable("$SJtest:RC$")
	_, err := v.Expand()
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestExpandStringVars_文字列中の変数を展開できる(t *testing.T) {
	AddSysValue("JOBNET", "ID", "123456")
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	res.St = "2015-01-01 12:34:56.789"
	res.Et = "2015-01-01 13:00:00.000"
	res.Var = "testout"
	AddJobValue("test", res)

	before := `NETWORK[$MSJOBNET:ID$] END. JOBRC[$MJtest:RC$] TESTENV[$METEST$]`
	after, err := ExpandStringVars(before, plcMaster, kndSys, kndEnv, kndJob)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	expect := `NETWORK[123456] END. JOBRC[1] TESTENV[testenv]`
	if after != expect {
		t.Errorf("変数展開後の文字列が想定と一致しない。")
		t.Logf("想定値：%s", expect)
		t.Logf("実績値：%s", after)
	}
}

func TestExpandStringVars_利用可能種別に指定した変数種別だけを展開する(t *testing.T) {
	AddSysValue("JOBNET", "ID", "123456")
	os.Setenv("TEST", "testenv")
	res := new(Response)
	res.RC = 1
	AddJobValue("test", res)

	before := `NETWORK[$MSJOBNET:ID$] END. JOBRC[$MJtest:RC$] TESTENV[$METEST$]`
	after, err := ExpandStringVars(before, plcMaster, kndEnv)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	expect := `NETWORK[$MSJOBNET:ID$] END. JOBRC[$MJtest:RC$] TESTENV[testenv]`
	if after != expect {
		t.Errorf("変数展開後の文字列が想定と一致しない。")
		t.Logf("想定値：%s", expect)
		t.Logf("実績値：%s", after)
	}
}

func TestExpandStringVars_場所識別子が異なる変数を無視する(t *testing.T) {
	os.Setenv("TEST", "testenv")

	before := `MASTERENV[$METEST$] SERVANTENV[$SETEST$]`
	after, err := ExpandStringVars(before, plcMaster, kndSys, kndEnv, kndJob)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	expect := `MASTERENV[testenv] SERVANTENV[$SETEST$]`
	if after != expect {
		t.Errorf("変数展開後の文字列が想定と一致しない。")
		t.Logf("想定値：%s", expect)
		t.Logf("実績値：%s", after)
	}
}

func TestExpandTime(t *testing.T) {
	v := NewVariable("$ST20150730123456.789$")
	result, err := v.expandTime()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}
	if result != "2015-07-30 21:34:56.789" {
		t.Errorf("変数展開後の文字列が想定と一致しない。")
		t.Logf("想定値：%s", "2015-07-30 21:34:56.789")
		t.Logf("実績値：%s", result)
	}
}
