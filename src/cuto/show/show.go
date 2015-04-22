// Copyright 2015 unirita Inc.
// Created 2015/04/14 shanxia

package main

import (
	"fmt"
	"os"

	"cuto/console"
	"cuto/db"
	"cuto/db/query"
	"cuto/show/gen"
)

// 表示に使用する構造体。
type ShowParam struct {
	jobnetName string         // ジョブネットワーク名
	from       string         // FROM日付
	to         string         // TO日付
	status     int            // ステータス
	gen        *gen.Generator // 出力ジェネレーター
	conn       db.IConnection // DBコネクション
}

// ジョブネットワークのインスタンス毎の表示用構造体
type oneJobnetwork struct {
	jobnet *db.JobNetworkResult // ジョブネットワーク情報
	jobs   []*db.JobResult      // ジョブネットワークに所属するジョブ情報一覧
}

// ShowParam構造体のコンストラクタ。
func NewShowParam(jobnetName string, from string, to string, status int, gen *gen.Generator) *ShowParam {
	return &ShowParam{
		jobnetName: jobnetName,
		from:       from,
		to:         to,
		status:     status,
		gen:        gen,
	}
}

// ユーティリティ実行のメインルーチン
// 成功した場合は、出力したジョブネットワークの件数を返します。
func (s *ShowParam) Run(db_name string) (int, error) {
	conn, err := db.Open(db_name)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	s.conn = conn

	// ジョブネットワーク情報の取得
	netResults, err := s.getJobnetworkList()
	if err != nil {
		return 0, err
	} else if len(netResults) == 0 {
		return 0, nil
	}
	// 取得したジョブネットワークインスタンス毎に、ジョブを出力する。
	var out gen.OutputRoot
	for _, jobnet := range netResults {
		oneJobnet := &oneJobnetwork{jobnet: jobnet}
		err := oneJobnet.getJobList(s.conn)
		if err != nil { // ジョブネットワーク内のジョブ取得に失敗したが、ジョブネットワークだけでも出力する。
			console.DisplayError("CTU005W", oneJobnet.jobnet.ID, err)
		}
		out.Jobnetworks = append(out.Jobnetworks, oneJobnet.setOutputStructure())
	}
	// ジェネレーターで出力メッセージ作成。
	msg, err := (*s.gen).Generate(&out)
	if err != nil {
		return 0, err
	}
	fmt.Fprint(os.Stdout, msg)
	return len(netResults), nil
}

// ジョブネットワーク一覧の取得
func (s *ShowParam) getJobnetworkList() ([]*db.JobNetworkResult, error) {
	jnQ := query.CreateJobnetworkQuery(s.conn)
	if len(s.jobnetName) > 0 {
		jnQ.AddAndWhereJobnetwork(s.jobnetName)
	}
	if len(s.from) > 0 {
		jnQ.AddAndWhereMoreThanStartdate(s.from)
	}
	if len(s.to) > 0 {
		jnQ.AddAndWhereLessThanStartdate(s.to)
	}
	if s.status != -1 {
		jnQ.AddAndWhereStatus(s.status)
	}
	jnQ.AddOrderBy(query.ORDERBY_ASC)
	netResults, err := jnQ.GetJobnetworkList()
	if err != nil {
		return nil, err
	}
	return netResults, nil
}

// ジョブネットワークに所属するジョブ情報一覧を取得
func (o *oneJobnetwork) getJobList(conn db.IConnection) error {
	var err error
	o.jobs, err = query.GetJobsOfTargetNetwork(conn, o.jobnet.ID, query.ORDERBY_ASC)
	if err != nil {
		return err
	}
	return nil
}

// 出力ジェネレータ構造体への格納
func (o *oneJobnetwork) setOutputStructure() *gen.OutputJobNet {
	jobNet := &gen.OutputJobNet{
		Id:         o.jobnet.ID,
		Jobnetwork: o.jobnet.JobnetWork,
		StartDate:  o.jobnet.StartDate,
		EndDate:    o.jobnet.EndDate,
		Status:     o.jobnet.Status,
		Detail:     o.jobnet.Detail,
		CreateDate: o.jobnet.CreateDate,
		UpdateDate: o.jobnet.UpdateDate,
	}
	for _, job := range o.jobs {
		j := &gen.OutputJob{
			JobId:      job.JobId,
			Jobname:    job.JobName,
			StartDate:  job.StartDate,
			EndDate:    job.EndDate,
			Status:     job.Status,
			Detail:     job.Detail,
			Rc:         job.Rc,
			Node:       job.Node,
			Port:       job.Port,
			Variable:   job.Variable,
			CreateDate: job.CreateDate,
			UpdateDate: job.UpdateDate,
		}
		jobNet.Jobs = append(jobNet.Jobs, j)
	}
	return jobNet
}
