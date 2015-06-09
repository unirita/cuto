// DB周りの定義や処理を記述するパッケージ。
package db

// 実行結果のステータス
const (
	RUNNING = iota
	NORMAL
	WARN
	ABNORMAL = 9
)

// 実行結果の文字列ステータス
const (
	ST_RUNNING  = "RUNNING"
	ST_NORMAL   = "NORMAL END"
	ST_WARN     = "WARN END"
	ST_ABNORMAL = "ABNORMAL END"
)
