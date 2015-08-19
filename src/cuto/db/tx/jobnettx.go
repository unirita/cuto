// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"fmt"

	"cuto/db"
	"cuto/db/query"
	"cuto/log"
	"cuto/utctime"
)

// ジョブIDをキーに持つ
type JobMap map[string]*db.JobResult

// ジョブ実行結果を保持する。
type ResultMap struct {
	JobnetResult *db.JobNetworkResult // ジョブネットワーク情報の構造体。
	Jobresults   JobMap               // ジョブネットワーク内のジョブ状態を保存するMap。
	conn         db.IConnection       // DBコネクション
}

// ジョブネットワークの開始状態を記録する。
//
// param : jobnetName ジョブネットワーク名。
//
// param : dbname データベース名。
//
// return : ジョブ実行結果を保持する構造体ポインタ。
//
// return : error
func StartJobNetwork(jobnetName string, dbname string) (*ResultMap, error) {
	jn := db.NewJobNetworkResult(jobnetName, utctime.Now().String(), db.RUNNING)

	conn, err := db.Open(dbname)
	if err != nil {
		return nil, err
	}
	resMap := &ResultMap{jn, make(JobMap), conn}

	if err := resMap.insertJobNetwork(); err != nil {
		return nil, err
	}
	return resMap, nil
}

// ジョブネットワークの実行結果を復元する。
func ResumeJobNetwork(nid int, dbname string) (*ResultMap, error) {
	conn, err := db.Open(dbname)
	if err != nil {
		return nil, err
	}
	jn, err := query.GetJobnetwork(conn, nid)
	if err != nil {
		return nil, err
	}

	jr, err := query.GetJobMapOfTargetNetwork(conn, nid)
	if err != nil {
		return nil, err
	}

	resMap := &ResultMap{
		JobnetResult: jn,
		Jobresults:   jr,
		conn:         conn,
	}

	return resMap, nil
}

// ネットワーク終了時に結果情報を設定する。同時にDBコネクションも切断する。
//
// param : status ジョブネットワークのステータス。
//
// param : detail ジョブネットワークに記録する詳細メッセージ。
//
// return : error
func (r *ResultMap) EndJobNetwork(status int, detail string) error {
	if r.conn == nil {
		return fmt.Errorf("Can't access DB file.")
	}
	defer r.conn.Close()

	if r.JobnetResult == nil {
		return fmt.Errorf("Invalid Jobnetwork info.")
	}
	r.JobnetResult.EndDate = utctime.Now().String()
	r.JobnetResult.Status = status
	r.JobnetResult.Detail = detail

	for _, jobresult := range r.Jobresults {
		if r.JobnetResult.Status < jobresult.Status {
			r.JobnetResult.Status = jobresult.Status
		}
	}

	if err := r.updateJobNetwork(); err != nil {
		return err
	}
	return nil
}

// DBコネクションを返す。
func (r *ResultMap) GetConnection() db.IConnection {
	return r.conn
}

// ジョブネットワークレコードをInsertする。
func (r *ResultMap) insertJobNetwork() error {
	var isCommit bool
	tx, err := r.conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()

	now := utctime.Now().String()
	r.JobnetResult.CreateDate = now
	r.JobnetResult.UpdateDate = now

	err = tx.Insert(r.JobnetResult)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("networkId[%v]", r.JobnetResult.ID))
	isCommit = true
	return nil
}

// ジョブネットワークレコードをUpdateする。
func (r *ResultMap) updateJobNetwork() error {
	var isCommit bool
	tx, err := r.conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()
	r.JobnetResult.UpdateDate = utctime.Now().String()

	if _, err = tx.Update(r.JobnetResult); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	isCommit = true
	return nil
}

// DBコネクションを外部から渡す。テスト用のメソッド。
func (r *ResultMap) SetConnection(conn db.IConnection) {
	r.conn = conn
}
