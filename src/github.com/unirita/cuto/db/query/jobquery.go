// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package query

import (
	"fmt"

	"github.com/unirita/cuto/db"
)

type jobQuery struct {
	sql  string         // SQL文
	conn db.IConnection // DBコネクション
}

// JOBテーブルの総件数を取得する。
func JobCountAll(conn db.IConnection) int {
	num, _ := conn.GetDbMap().SelectInt("select count(*) from JOB")
	return int(num)
}

//　ジョブネットワークのインスタンスIDを指定して、ジョブ情報を取得する。
//
// param - conn 接続済みのDBコネクション。
//
// param - nid 検索に使用するジョブネットワークのインスタンスID。
//
// param - orderby 昇順（query.ORDERBY_ASC) / 降順(query.ORDERBY_DESC)
//
// return ジョブネットワークレコードのスライスとエラー情報
func GetJobsOfTargetNetwork(conn db.IConnection, nid int, orderby int) ([]*db.JobResult, error) {
	q := CreateJobQuery(conn)
	q.AddAndWhereID(nid)
	q.AddOrderBy(orderby)

	list, err := conn.GetDbMap().Select(db.JobResult{}, q.sql)
	if err != nil {
		return nil, err
	}
	var results []*db.JobResult
	for _, l := range list {
		r := l.(*db.JobResult)
		results = append(results, r)
	}
	return results, nil
}

//　ジョブネットワークのインスタンスIDを指定して、ジョブ情報をマップ形式で取得する。
func GetJobMapOfTargetNetwork(conn db.IConnection, nid int) (map[string]*db.JobResult, error) {
	q := CreateJobQuery(conn)
	q.AddAndWhereID(nid)

	list, err := conn.GetDbMap().Select(db.JobResult{}, q.sql)
	if err != nil {
		return nil, err
	}

	results := make(map[string]*db.JobResult)
	for _, l := range list {
		r := l.(*db.JobResult)
		results[r.JobId] = r
	}
	return results, nil
}

func CreateJobQuery(conn db.IConnection) *jobQuery {
	sql := fmt.Sprintf("select ID,JOBID,JOBNAME,STARTDATE,ENDDATE,STATUS,DETAIL,RC,NODE,PORT,VARIABLE,CREATEDATE,UPDATEDATE from JOB where 0=0 ")
	return &jobQuery{sql, conn}
}

// 引数に指定したIDと合致する条件を追加。
func (j *jobQuery) AddAndWhereID(id int) {
	j.sql = fmt.Sprintf(" %v and ID = %v ", j.sql, id)
}

// ORDER BY句を追加する。
// 引数へは ORDERBY_ASC または ORDERBY_DESC を指定する。
func (j *jobQuery) AddOrderBy(orderby int) {
	if orderby == ORDERBY_ASC {
		j.sql = fmt.Sprintf("%s order by UPDATEDATE asc ", j.sql)
	} else {
		j.sql = fmt.Sprintf("%s order by UPDATEDATE desc ", j.sql)
	}
}
