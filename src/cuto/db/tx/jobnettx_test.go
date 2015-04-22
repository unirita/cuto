// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"fmt"
	"os"
	"testing"

	"cuto/db"
	"cuto/db/query"
)

// テストDB名
var db_path = fmt.Sprintf("%s%c%s%c%s%c%s%c%s", os.Getenv("GOPATH"), os.PathSeparator, "test",
	os.PathSeparator, "cuto", os.PathSeparator, "db", os.PathSeparator, "tx")
var db_name = fmt.Sprintf("%s%c%s", db_path, os.PathSeparator, "test_tx.sqlite")

// DB接続後の失敗を誘うためのダミーファイル
var dummy_db = fmt.Sprintf("%v%c%v", db_path, os.PathSeparator, "dummy.sqlite")

// DBの初期化。
func init() {
	conn, err := db.Open(db_name)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	conn.GetDbMap().DropTables()
	err = conn.GetDbMap().CreateTables()
	if err != nil {
		panic(err.Error())
	}
	// テストを繰り返すとDBが肥大化する対策
	conn.GetDb().Exec("vacuum")
}

// DBを検査して、登録件数と内容を取得する。
func verifyDb(nid int) (int, *db.JobNetworkResult) {
	conn, err := db.Open(db_name)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	num := query.JobnetworkCountAll(conn)

	res, err := query.GetJobnetwork(conn, nid)
	if err != nil {
		panic(err)
	}
	return num, res
}

func TestStartJobNetwork_ジョブネットワークの開始ステータスが正常か(t *testing.T) {
	name := "JNet1"

	resMap, err := StartJobNetwork(name, db_name)
	if err != nil {
		t.Errorf("エラーがすべきでないパターンで、エラーが発生しました。: %s", err.Error())
	}
	if resMap.JobnetResult.Status != db.RUNNING {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]になっています。", db.RUNNING, resMap.JobnetResult.Status)
	}
	// DB検証
	num, res := verifyDb(resMap.JobnetResult.ID)
	if num != 1 {
		t.Errorf("登録件数が1件のはずが、[%v]件見つかりました。", num)
	}
	if res.JobnetWork != name {
		t.Errorf("[%v]が見つかるはずが、異なるジョブネット[%v]が返りました。", name, res.JobnetWork)
	}

	c := resMap.GetConnection()
	if c == nil {
		t.Error("コネクションの取得に失敗しました。")
	}
}

func TestStartJobNetwork_DBOpenに失敗(t *testing.T) {
	name := "JNet1"

	_, err := StartJobNetwork(name, "db_name")
	if err == nil {
		t.Error("存在しないDBファイルを指定したのに、エラーが返りませんでした")
	}
}

func TestStartJobNetwork_Insertに失敗(t *testing.T) {
	name := "JNet1"

	dummy_db := fmt.Sprintf("%v%c%v", db_path, os.PathSeparator, "dummy.sqlite")

	_, err := StartJobNetwork(name, dummy_db)
	if err == nil {
		t.Error("存在しないDBファイルを指定したのに、エラーが返りませんでした")
	}
}

func TestEndJobNetwork_ジョブネットワークの開始ステータスに正常を登録(t *testing.T) {
	// ここからはテスト前データ登録 ///////
	name := "JNet2"

	resMap, err := StartJobNetwork(name, db_name)
	if resMap.JobnetResult.Status != db.RUNNING {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]になっています。", db.RUNNING, resMap.JobnetResult.Status)
	}
	if err != nil {
		t.Errorf("エラーがすべきでないパターンで、エラーが発生しました。: %s", err.Error())
	}
	// ここまで ///////

	resMap.EndJobNetwork(db.NORMAL, "")

	if resMap.JobnetResult.Status != db.NORMAL {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]が返りました。", db.NORMAL, resMap.JobnetResult.Status)
	}
	// DB検証
	num, res := verifyDb(resMap.JobnetResult.ID)
	if num != 2 {
		t.Errorf("登録件数が2件のはずが、[%v]件見つかりました。", num)
	}
	if res.JobnetWork != name {
		t.Errorf("[%v]が見つかるはずが、異なるジョブネット[%v]が返りました。", name, res.JobnetWork)
	}
	if res.Status != db.NORMAL {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]がDBに登録されていました。", db.NORMAL, res.Status)
	}
}

func TestEndJobNetwork_ジョブネットワークの開始ステータスに異常を登録(t *testing.T) {
	// ここからはテスト前データ登録 ///////
	name := "JNet3"

	resMap, err := StartJobNetwork(name, db_name)
	if resMap.JobnetResult.Status != db.RUNNING {
		t.Fatalf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]になっています。", db.RUNNING, resMap.JobnetResult.Status)
	}
	if err != nil {
		t.Fatalf("エラーがすべきでないパターンで、エラーが発生しました。: %s", err.Error())
	}
	// ここまで ///////
	detail := "An accident occurred."

	resMap.EndJobNetwork(db.ABNORMAL, detail)
	if resMap.JobnetResult.Status != db.ABNORMAL {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]が返りました。", db.ABNORMAL, resMap.JobnetResult.Status)
	}
	if resMap.JobnetResult.Detail != detail {
		t.Errorf("ジョブネットワークの詳細メッセージ[%v]が返るべきところ、[%v]が返りました。", detail, resMap.JobnetResult.Detail)
	}
	// DB検証
	num, res := verifyDb(resMap.JobnetResult.ID)
	if num != 3 {
		t.Errorf("登録件数が3件のはずが、[%v]件見つかりました。", num)
	}
	if res.JobnetWork != name {
		t.Errorf("[%v]が見つかるはずが、異なるジョブネット[%v]が返りました。", name, res.JobnetWork)
	}
	if res.Status != db.ABNORMAL {
		t.Errorf("ジョブネットワークのステータスが[%v]になるべきところ、[%v]がDBに登録されていました。", db.NORMAL, res.Status)
	}
}

func TestEndJobNetwork_不正なジョブネットワーク情報(t *testing.T) {
	con, err := db.Open(dummy_db)
	if err != nil {
		t.Fatalf("DB接続に失敗しました。 - %v", err)
	}
	resMap := &ResultMap{
		conn: con,
	}
	err = resMap.EndJobNetwork(0, "")
	if err == nil {
		t.Error("エラーが返るべきところ、成功しました。")
	}
}

func TestEndJobNetwork_更新失敗(t *testing.T) {
	resMap := &ResultMap{}
	err := resMap.EndJobNetwork(0, "")
	if err == nil {
		t.Error("エラーが返るべきところ、成功しました。")
	}
}

func TestEndJobNetwork_接続失敗(t *testing.T) {
	con, err := db.Open(dummy_db)
	if err != nil {
		t.Fatalf("DB接続に失敗しました。 - %v", err)
	}
	jobnet := &db.JobNetworkResult{}
	resMap := &ResultMap{
		JobnetResult: jobnet,
		conn:         con,
	}
	err = resMap.EndJobNetwork(0, "")
	if err == nil {
		t.Error("エラーが返るべきところ、成功しました。")
	}
}
