// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"sync"
	"testing"

	"cuto/db"
	"cuto/utctime"
)

var mutex sync.Mutex

func TestInsertJob_ジョブの新規登録処理(t *testing.T) {
	conn, err := db.Open(db_name)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = InsertJob(conn,
		&db.JobResult{
			ID:        121,
			JobId:     "JOB001",
			JobName:   "abcjob.bat",
			StartDate: utctime.Now().String(),
			Status:    db.RUNNING,
			Node:      "TestNode01",
			Port:      9999,
		}, &mutex)
	if err != nil {
		t.Error("ジョブテーブルへの登録に失敗しました。 - ", err)
	}
}

func TestInsertJob_ジョブの新規登録失敗(t *testing.T) {
	conn, err := db.Open(dummy_db)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = InsertJob(conn,
		&db.JobResult{
			ID:        131,
			JobId:     "JOB001",
			JobName:   "abcjob.bat",
			StartDate: utctime.Now().String(),
			Status:    db.RUNNING,
			Node:      "TestNode01",
			Port:      9999,
		}, &mutex)
	if err == nil {
		t.Error("予定していた失敗が返りませんでした。 - ")
	}
}

func TestUpdateJob_ジョブの更新処理(t *testing.T) {
	conn, err := db.Open(db_name)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	jobres := &db.JobResult{
		ID:        122,
		JobId:     "JOB002",
		JobName:   "XYZ.vbs",
		StartDate: utctime.Now().String(),
		Status:    db.RUNNING,
		Node:      "TestNode02",
		Port:      9999,
	}
	err = InsertJob(conn, jobres, &mutex)
	if err != nil {
		t.Fatal("ジョブテーブルへの登録に失敗しました。 - ", err)
	}
	jobres.Status = db.ABNORMAL
	jobres.Variable = "VAR"
	jobres.EndDate = utctime.Now().String()
	jobres.Detail = "Attention!!"
	jobres.Rc = 4

	err = UpdateJob(conn, jobres, &mutex)
	if err != nil {
		t.Error("ジョブテーブルの更新に失敗しました。 - ", err)
	}
}

func TestUpdateJob_ジョブの登録前に更新(t *testing.T) {
	conn, err := db.Open(dummy_db)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	jobres := &db.JobResult{
		ID:        132,
		JobId:     "JOB002",
		JobName:   "XYZ.vbs",
		StartDate: utctime.Now().String(),
		Status:    db.RUNNING,
		Node:      "TestNode02",
		Port:      9999,
	}
	err = UpdateJob(conn, jobres, &mutex)
	if err == nil {
		t.Error("予定していた失敗が返りませんでした。 - ")
	}
}

func TestUpdateJob_ジョブの更新失敗(t *testing.T) {
	conn, err := db.Open(dummy_db)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	jobres := &db.JobResult{
		ID:        122,
		JobId:     "JOB002",
		JobName:   "XYZ.vbs",
		StartDate: utctime.Now().String(),
		Status:    db.RUNNING,
		Node:      "TestNode02",
		Port:      9999,
	}
	err = UpdateJob(conn, jobres, &mutex)
	if err == nil {
		t.Error("予定していた失敗が返りませんでした。 - ")
	}
}
