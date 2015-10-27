// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package query

import (
	"testing"
)

const all_jobnum = 12 // JOBテーブルに存在する総件数。

func TestJobCountAll_件数取得(t *testing.T) {
	num := JobCountAll(conn)
	if num != all_jobnum {
		t.Error("テストデータが%v件のはずが、[%v]件が返ってきました。", all_jobnum, num)
	}
}

func TestGetJobsOfTargetNetwork_ジョブネットワークIDを指定して取得(t *testing.T) {
	results, err := GetJobsOfTargetNetwork(conn, 4, ORDERBY_ASC)
	if err != nil {
		t.Error("エラーが返ってきました。 - ", err)
	}
	if len(results) != 2 {
		t.Errorf("2件見つかるべきところ、%v件が返りました。", len(results))
	}
	if results[0].JobId != "最適化" {
		t.Errorf("不正なジョブID[%v]が返りました。", results[0].JobId)
	}
	if results[1].JobId != "バックアップ" {
		t.Errorf("不正なジョブID[%v]が返りました。", results[1].JobId)
	}
}

func TestGetJobsOfTargetNetwork_ジョブを0件取得(t *testing.T) {
	results, err := GetJobsOfTargetNetwork(conn, 999, ORDERBY_DESC)
	if err != nil {
		t.Error("エラーが返ってきました。 - ", err)
	}
	if len(results) != 0 {
		t.Errorf("0件が返るべきところ、%v件が返ってきました。 - ", len(results))
	}
}
