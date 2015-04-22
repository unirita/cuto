// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package query

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"cuto/db"
)

// テストDB名
var (
	db_root = fmt.Sprintf("%s%c%s%c%s%c%s%c%s", os.Getenv("GOPATH"), os.PathSeparator, "test", os.PathSeparator, "cuto", os.PathSeparator, "db", os.PathSeparator, "query")
	db_name = fmt.Sprintf("%s%c%s", db_root, os.PathSeparator, "test_q.sqlite")
	conn    = dbOpen()

	all_jobnetNum = 6
)

func dbOpen() db.IConnection {
	c, err := db.Open(db_name)
	if err != nil {
		panic(err)
	}
	return c
}

func TestJobnetworkCountAll_件数取得(t *testing.T) {
	num := JobnetworkCountAll(conn)
	if num != all_jobnetNum {
		t.Errorf("テストデータが%v件のはずが、[%v]件が返ってきました。", all_jobnetNum, num)
	}
}

func TestGetJobnetwork_1件取得(t *testing.T) {
	result, err := GetJobnetwork(conn, 2)
	if err != nil {
		t.Error("ジョブ取得時にエラーが返ってきました。 - ", err)
	}
	if result.ID != 2 {
		t.Errorf("指定と異なるジョブネットID[%v]が返ってきました。", result.ID)
	}
	if result.JobnetWork != "ジョブネット2" {
		t.Errorf("指定と異なるジョブネット[%v]が返ってきました。", result.JobnetWork)
	}
	if result.Status != 0 {
		t.Errorf("指定と異なるジョブネットのステータス[%v]が返ってきました。", result.Status)
	}
}

func TestGetJobnetwork_0件取得(t *testing.T) {
	result, err := GetJobnetwork(conn, 999)
	if err == nil {
		t.Error("ジョブ取得時に失敗しても、エラーが返ってきませんでした。")
	}
	if result != nil {
		t.Error("ジョブが0件なのにnilが返りませんでした。")
	}
}

func TestGetJobnetwork_クエリ発行に失敗(t *testing.T) {
	_, err := GetJobnetwork(nil, 999)
	if err == nil {
		t.Error("ダミーのコネクションを渡しても、エラーが返ってきませんでした。")
	}
}

func TestGetJobnetworkListFromName_41件取得(t *testing.T) {
	const jobnet string = "ジョブネット1"
	const count int = 2

	results, err := GetJobnetworkListFromName(conn, jobnet, ORDERBY_DESC)
	if err != nil {
		t.Error("ジョブ取得時にエラーが返りました。 - ", err)
	}
	if len(results) != count {
		t.Errorf("ジョブが%v件の想定ですが、返ってきたのは[%v]件でした。", count, len(results))
	}
	if results[0].ID != 5 {
		t.Errorf("1件目のジョブネットIDは5の想定ですが、%vが返りました。", results[0].ID)
	}
	if results[1].ID != 1 {
		t.Errorf("1件目のジョブネットIDは5の想定ですが、%vが返りました。", results[0].ID)
	}
}

func TestGetJobnetworkList_2件取得(t *testing.T) {
	query := CreateJobnetworkQuery(conn)
	query.AddAndWhereJobnetwork("ジョブネット1")
	query.AddOrderBy(ORDERBY_ASC)
	if !strings.Contains(query.sql, "asc") {
		t.Errorf("sqlが不正です。 - %v", query.sql)
	}
	results, err := query.GetJobnetworkList()
	if err != nil {
		t.Error("ジョブ取得時にエラーが返ってきました。 - ", err)
	}
	if len(results) != 2 {
		t.Errorf("2件返ってくるべきところ、%v件が返ってきました。", len(results))
	}
	if results[1].ID != 5 {
		t.Errorf("5が返ってくるべきところ、%vが返ってきました。", results[0].ID)
	}
	if results[0].ID != 1 {
		t.Errorf("1が返ってくるべきところ、%vが返ってきました。", results[1].ID)
	}
}

func TestGetJobnetworkList_クエリ不正(t *testing.T) {
	query := CreateJobnetworkQuery(conn)
	query.AddAndWhereJobnetwork("ジョブネット1")
	query.AddOrderBy(ORDERBY_DESC)
	if !strings.Contains(query.sql, "desc") {
		t.Errorf("sqlが不正です。 - %v", query.sql)
	}
	query.sql = "abc" // 不正なクエリに書き換え
	_, err := query.GetJobnetworkList()
	if err == nil {
		t.Error("不正なクエリを発行したのに、失敗しませんでした。")
	}
}

func TestAddAndWhereLessThanStartdate_開始日より過去を取得(t *testing.T) {
	query := CreateJobnetworkQuery(conn)
	query.AddAndWhereLessThanStartdate("2015-04-17 09:08:07.000")
	if !strings.Contains(query.sql, "STARTDATE < '2015-04-17 09:08:07.000'") {
		t.Errorf("不正なSQLです。 - %v", query.sql)
	}
}

func TestAddAndWhereMoreThanStartdate_開始日より過去を取得(t *testing.T) {
	query := CreateJobnetworkQuery(conn)
	query.AddAndWhereMoreThanStartdate("2015-04-17 09:08:07.000")
	if !strings.Contains(query.sql, "'2015-04-17 09:08:07.000' < STARTDATE") {
		t.Errorf("不正なSQLです。 - %v", query.sql)
	}
}

func TestAddAndWhereStatus_Statusでフィルタ(t *testing.T) {
	query := CreateJobnetworkQuery(conn)
	query.AddAndWhereStatus(9)
	if !strings.Contains(query.sql, "STATUS = 9") {
		t.Errorf("不正なSQLです。 - %v", query.sql)
	}
}
