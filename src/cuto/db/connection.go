// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

import (
	"errors"
	"fmt"
	"os"

	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type IConnection interface {
	GetDbMap() *gorp.DbMap
	GetDb() *sql.DB
	Close()
}

type Connection struct {
	db    *sql.DB     // DBコネクション
	dbMap *gorp.DbMap // マッピング情報
}

var sqlite3_driver = "sqlite3"

// SQLite3のセッションを接続し、テーブルとDTOのマッピングを行う。
//
// param - dbfile sqliteファイルのパス。
//
// return - コネクション情報とエラー情報
func Open(dbfile string) (IConnection, error) {
	if _, exist := os.Stat(dbfile); exist != nil {
		return nil, errors.New(fmt.Sprintf("Not found dbfile[%v]", dbfile))
	}
	db, err := sql.Open(sqlite3_driver, dbfile)
	if err != nil {
		return nil, err
	}
	// 外部キーを有効にする。
	db.Exec("PRAGMA foreign_keys=ON;")

	// テーブルと構造体のマッピング。
	dbmap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.SqliteDialect{},
	}
	jobNetworkMapping(dbmap)
	jobMapping(dbmap)

	return Connection{db, dbmap}, nil
}

func jobNetworkMapping(dbmap *gorp.DbMap) {
	t := dbmap.AddTableWithName(JobNetworkResult{}, "JOBNETWORK").SetKeys(true, "ID")
	t.ColMap("JobnetWork").Rename("JOBNETWORK")
	t.ColMap("StartDate").Rename("STARTDATE")
	t.ColMap("EndDate").Rename("ENDDATE")
	t.ColMap("Status").Rename("STATUS")
	t.ColMap("Detail").Rename("DETAIL")
	t.ColMap("CreateDate").Rename("CREATEDATE")
	t.ColMap("UpdateDate").Rename("UPDATEDATE")
}

func jobMapping(dbmap *gorp.DbMap) {
	t := dbmap.AddTableWithName(JobResult{}, "JOB").SetKeys(false, "ID", "JobId")
	t.ColMap("JobId").Rename("JOBID")
	t.ColMap("JobName").Rename("JOBNAME")
	t.ColMap("StartDate").Rename("STARTDATE")
	t.ColMap("EndDate").Rename("ENDDATE")
	t.ColMap("Status").Rename("STATUS")
	t.ColMap("Detail").Rename("DETAIL")
	t.ColMap("Rc").Rename("RC")
	t.ColMap("Node").Rename("NODE")
	t.ColMap("Port").Rename("PORT")
	t.ColMap("Variable").Rename("VARIABLE")
	t.ColMap("CreateDate").Rename("CREATEDATE")
	t.ColMap("UpdateDate").Rename("UPDATEDATE")
}

// SQLite3とのセッションを切断する。
func (c Connection) Close() {
	c.db.Close()
}

// DBマッピング情報を返す。
func (c Connection) GetDbMap() *gorp.DbMap {
	return c.dbMap
}

// DBオブジェクトを返す。
func (c Connection) GetDb() *sql.DB {
	return c.db
}
