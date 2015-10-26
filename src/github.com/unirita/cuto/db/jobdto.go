// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

// ジョブ実行結果
type JobResult struct {
	ID         int    // ジョブネットワークのインシデントID
	JobId      string // ジョブID
	JobName    string // ジョブ名
	StartDate  string // ジョブの起動日時
	EndDate    string // ジョブの終了日時
	Status     int    // ステータス
	Detail     string // 詳細メッセージ
	Rc         int    // リターンコード
	Node       string // ノード名
	Port       int    // ポート番号
	Variable   string // 変数情報
	CreateDate string // 作成日時
	UpdateDate string // 更新日時
}

// ジョブ実行結果のコンストラクタ。
//
// param : id ジョブネットワークのインシデントID
//
// return : JobResultポインタ
func NewJobResult(id int) *JobResult {
	return &JobResult{ID: id}
}
