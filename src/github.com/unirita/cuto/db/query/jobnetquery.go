// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package query

import (
	"fmt"

	"github.com/unirita/cuto/db"
)

const (
	ORDERBY_ASC = iota
	ORDERBY_DESC
)

type JobNetResultQuery struct {
	sql  string         // SQL文
	conn db.IConnection // コネクション
}

// JOBNETWORKテーブルの総件数を取得する。
func JobnetworkCountAll(conn db.IConnection) int {
	num, _ := conn.GetDbMap().SelectInt("select count(*) from JOBNETWORK")
	return int(num)
}

func CreateJobnetworkQuery(conn db.IConnection) *JobNetResultQuery {
	sql := fmt.Sprintf("select ID,JOBNETWORK,STARTDATE,ENDDATE,STATUS,DETAIL,PID,CREATEDATE,UPDATEDATE from JOBNETWORK where 0=0 ")
	return &JobNetResultQuery{sql, conn}
}

//　ジョブネットワークのインスタンスIDを指定して、ジョブネットワーク詳細情報を取得する。
func GetJobnetwork(conn db.IConnection, id int) (*db.JobNetworkResult, error) {
	q := CreateJobnetworkQuery(conn)
	q.AddAndWhereID(id)

	results, err := q.GetJobnetworkList()
	if err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, fmt.Errorf("Network[id = %d] not found.", id)
	}
	return results[0], nil
}

// ジョブネットワーク名を指定して、一覧を取得する。
//
// param - conn 接続済みのDBコネクション。
//
// param - name 検索に使用するジョブネットワーク名。
//
// param - orderby 昇順（query.ORDERBY_ASC) / 降順(query.ORDERBY_DESC)
//
// return ジョブネットワークレコードのスライスとエラー情報
func GetJobnetworkListFromName(conn db.IConnection, name string, orderby int) ([]*db.JobNetworkResult, error) {
	q := CreateJobnetworkQuery(conn)
	q.AddAndWhereJobnetwork(name)
	q.AddOrderBy(orderby)

	return q.GetJobnetworkList()
}

// ジョブネットワーク一覧を取得する。
func (j *JobNetResultQuery) GetJobnetworkList() ([]*db.JobNetworkResult, error) {
	if j.conn == nil {
		return nil, fmt.Errorf("Invalid DB Connection.")
	}
	list, err := j.conn.GetDbMap().Select(db.JobNetworkResult{}, j.sql)
	if err != nil {
		return nil, err
	}
	var results []*db.JobNetworkResult
	for _, l := range list {
		r := l.(*db.JobNetworkResult)
		results = append(results, r)
	}
	return results, nil
}

// 引数に指定したIDと合致する条件を追加。
func (j *JobNetResultQuery) AddAndWhereID(id int) {
	j.sql = fmt.Sprintf(" %v and ID = %v ", j.sql, id)
}

// 引数に指定したJOBNETWORKと合致する条件を追加。
func (j *JobNetResultQuery) AddAndWhereJobnetwork(jobnetwork string) {
	j.sql = fmt.Sprintf(" %v and JOBNETWORK = '%v' ", j.sql, jobnetwork)
}

// 引数に指定したSTARTDATEよりも小さい日付[ STARTDATE < '引数' ]を取得。
func (j *JobNetResultQuery) AddAndWhereLessThanStartdate(startDate string) {
	j.sql = fmt.Sprintf(" %v and STARTDATE < '%v' ", j.sql, startDate)
}

// 引数に指定したSTARTDATEよりも大きい日付[ '引数' < STARTDATE ]を取得。
func (j *JobNetResultQuery) AddAndWhereMoreThanStartdate(startDate string) {
	j.sql = fmt.Sprintf(" %v and '%v' < STARTDATE ", j.sql, startDate)
}

// 引数に指定したSTATUSと合致する条件を追加。
func (j *JobNetResultQuery) AddAndWhereStatus(status int) {
	j.sql = fmt.Sprintf(" %v and STATUS = %v ", j.sql, status)
}

// ORDER BY句を追加する。
// 引数へは ORDERBY_ASC または ORDERBY_DESC を指定する。
func (j *JobNetResultQuery) AddOrderBy(orderby int) {
	if orderby == ORDERBY_ASC {
		j.sql = fmt.Sprintf("%s order by UPDATEDATE asc ", j.sql)
	} else {
		j.sql = fmt.Sprintf("%s order by UPDATEDATE desc ", j.sql)
	}
}
