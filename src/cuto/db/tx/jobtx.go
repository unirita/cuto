// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"sync"
	"time"

	"cuto/db"
	"cuto/util"
)

// JOBテーブルへINSERTする。
//
// param - conn DBコネクション
//
// param - job JOBレコード構造体ポインタ
func InsertJob(conn db.IConnection, job *db.JobResult, mutex *sync.Mutex) error {
	mutex.Lock()
	defer mutex.Unlock()

	var isCommit bool

	tx, err := conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()

	now := util.DateFormat(time.Now())
	job.CreateDate = now
	job.UpdateDate = now

	err = tx.Insert(job)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	isCommit = true
	return nil
}

// JOBテーブルへUPDATEする。
//
// param - conn DBコネクション
//
// param - job JOBレコード構造体ポインタ
func UpdateJob(conn db.IConnection, job *db.JobResult, mutex *sync.Mutex) error {
	mutex.Lock()
	defer mutex.Unlock()

	var isCommit bool
	tx, err := conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()

	job.UpdateDate = util.DateFormat(time.Now())

	if _, err = tx.Update(job); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	isCommit = true
	return nil
}
