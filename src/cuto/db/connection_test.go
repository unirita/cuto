// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

import (
	"fmt"
	"os"

	"testing"
)

var (
	dbfile = getDBFile()
)

func getDBFile() string {
	return fmt.Sprintf("%v%c%v%c%v%c%v%c%v", os.Getenv("GOPATH"),
		os.PathSeparator, "test", os.PathSeparator, "cuto", os.PathSeparator, "db", os.PathSeparator, "test.sqlite")
}

func TestOpen_DBコネクションがOpenできる(t *testing.T) {
	con, err := Open(dbfile)
	if err != nil {
		t.Errorf("DBとの接続に失敗しました。 - %v", err)
	} else if con == nil {
		t.Error("Connectionオブジェクトがnilを返しました。")
	}
	defer con.Close()
}

func TestOpen_存在しないファイルを指定する(t *testing.T) {
	testfile := "xxx.testDB"
	con, err := Open(testfile)
	if err == nil {
		t.Errorf("DBとの接続に失敗しなければならないところ、成功が返りました。")
	} else if con != nil {
		defer con.Close()
	}
}

func TestOpen_存在しないドライバを指定する(t *testing.T) {
	bk := sqlite3_driver
	sqlite3_driver = "sqlite2" // 存在しないドライバ名
	defer func(driver string) {
		sqlite3_driver = driver
	}(bk)

	con, err := Open(dbfile) // 存在しないDBファイルを指定しても成功してしまう。
	if err == nil {
		t.Errorf("エラーを返すべきところ、エラー情報がnilを返しました。")
	}
	if con != nil {
		t.Error("エラーが返ってきたのに、connectionオブジェクトがnilではなかった。")
	}
}

func TestGetDbMap_DBMapを取得する(t *testing.T) {
	con, err := Open(dbfile)
	if err != nil {
		t.Fatalf("DBとの接続に失敗しました。 - %v", err)
	} else if con == nil {
		t.Fatal("Connectionオブジェクトがnilを返しました。")
	}
	defer con.Close()

	if dbMap := con.GetDbMap(); dbMap == nil {
		t.Error("DBMapが取得できません。")
	}
}

func TestGetDb_DBオブジェクトを取得する(t *testing.T) {
	con, err := Open(dbfile)
	if err != nil {
		t.Fatalf("DBとの接続に失敗しました。 - %v", err)
	} else if con == nil {
		t.Fatal("Connectionオブジェクトがnilを返しました。")
	}
	defer con.Close()

	if db := con.GetDb(); db == nil {
		t.Error("DBオブジェクトが取得できません。")
	}
}
