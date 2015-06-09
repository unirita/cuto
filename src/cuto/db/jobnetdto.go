// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

// ジョブネットワーク実行結果
type JobNetworkResult struct {
	ID         int    // ジョブネットワークのインシデントID
	JobnetWork string // ジョブネットワーク名
	StartDate  string // 起動日時
	EndDate    string // 終了日時
	Status     int    // ジョブネットワークのステータス
	Detail     string // 詳細メッセージ
	CreateDate string // 作成日時
	UpdateDate string // 更新日時
}

// ジョブネットワーク実行結果のコンストラクタ。
//
// param : jobnetName ジョブネットワーク名。
//
// param : startDate ジョブネットワーク起動日時。
//
// param : status ステータス。
//
// return : JobNetworkResultポインタ
func NewJobNetworkResult(jobnetName string, startDate string, status int) *JobNetworkResult {
	return &JobNetworkResult{
		JobnetWork: jobnetName,
		StartDate:  startDate,
		Status:     status,
	}
}
