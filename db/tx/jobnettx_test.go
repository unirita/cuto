// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"path/filepath"
	"testing"

	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/db/query"
	"github.com/unirita/cuto/testutil"
)

// テストDB名
var db_path = filepath.Join(testutil.GetBaseDir(), "db", "tx", "_testdata")
var db_name = filepath.Join(db_path, "test_tx.sqlite")

// DB接続後の失敗を誘うためのダミーファイル
var dummy_db = filepath.Join(db_path, "dummy.sqlite")

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

	_, err := StartJobNetwork(name, dummy_db)
	if err == nil {
		t.Error("存在しないDBファイルを指定したのに、エラーが返りませんでした")
	}
}

func TestResumeJobNetwork_前回の実行実績を取得できる(t *testing.T) {
	dbfile := filepath.Join(db_path, "resume.sqlite")
	res, err := ResumeJobNetwork(1, dbfile)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	if res.JobnetResult.JobnetWork != "JNet1" {
		t.Errorf("ジョブネット[%s]が見つかるはずが、異なるジョブネット[%s]が返りました。", res.JobnetResult.JobnetWork, "JNet1")
	}
	if len(res.jobresults) != 2 {
		t.Fatalf("ジョブ実行結果の件数が%d件取得されるはずが、%d件取得された。", 2, len(res.jobresults))
	}
	if _, ok := res.jobresults["JOB001"]; !ok {
		t.Errorf("ジョブ[%s]の実行結果が取得されるはずが、されなかった。", "JOB001")
	}
	if _, ok := res.jobresults["JOB002"]; !ok {
		t.Errorf("ジョブ[%s]の実行結果が取得されるはずが、されなかった。", "JOB002")
	}
}

func TestResumeJobNetwork_DBOpenに失敗(t *testing.T) {
	_, err := ResumeJobNetwork(1, "db_name")
	if err == nil {
		t.Error("存在しないDBファイルを指定したのに、エラーが返りませんでした")
	}
}

func TestResumeJobNetwork_Selectに失敗(t *testing.T) {
	_, err := ResumeJobNetwork(1, dummy_db)
	if err == nil {
		t.Error("不正なDBファイルを指定したのに、エラーが返りませんでした")
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
