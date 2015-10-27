// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package parser

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/unirita/cuto/log"
)

// 拡張ジョブ情報
type JobEx struct {
	Node          string // ノード名
	Port          int    // ポート番号
	FilePath      string // ジョブファイルパス
	Param         string // 実行時引数
	Env           string // 実行時環境変数
	Workspace     string // 作業フォルダ
	WrnRC         int    // 警告終了判断に使用するリターンコードの下限値
	WrnPtn        string // 警告終了と判断する出力文字列
	ErrRC         int    // 異常終了判断に使用するリターンコードの下限値
	ErrPtn        string // 異常終了と判断する出力文字列
	TimeoutMin    int    // タイムアウト（分）
	SecondaryNode string // ノード名
	SecondaryPort int    // ポート番号
}

// CSVファイルの項目数
const (
	noSecondary   = 12
	withSecondary = 14
)

// 項目のインデックス
const (
	nameIdx = iota
	nodeIdx
	portIdx
	pathIdx
	paramIdx
	envIdx
	workIdx
	wrcIdx
	wptIdx
	ercIdx
	eptIdx
	tmoutIdx
	secNodeIdx
	secPortIdx
)

// JobEx構造体のオブジェクトを生成しする。
//
// return : 生成したオブジェクト
func NewJobEx() *JobEx {
	je := new(JobEx)
	je.TimeoutMin = -1
	return je
}

// ファイルから拡張ジョブ定義CSVを読み込み、パースする。
//
// param : fileName　ファイル名。
//
// return : 拡張ジョブ情報のパース後Map。
//
// return : エラー情報。
func ParseJobExFile(fileName string) (map[string]*JobEx, error) {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// ファイルが存在しない場合は空のマップを返す
			return make(map[string]*JobEx), nil
		} else {
			return nil, err
		}
	}
	defer file.Close()

	return ParseJobEx(file)
}

// readerから読み込んだ拡張ジョブ定義CSVをパースする。
// 空のカラムにはゼロ値をセットする。
//
// param : reader ファイルリーダー。
//
// return : 拡張ジョブ情報のパース後Map。
//
// return : エラー情報。
func ParseJobEx(reader io.Reader) (map[string]*JobEx, error) {
	r := csv.NewReader(reader)
	jobExMap := make(map[string]*JobEx)

	for i := 1; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if i == 1 {
			// タイトル行を無視する
			continue
		}

		if len(record) != noSecondary && len(record) != withSecondary {
			log.Info("Jobex line[%d] was ignored: Irregal column count[%d].", i, len(record))
			continue
		}

		name := record[nameIdx]
		if len(name) == 0 {
			log.Info("Jobex line[%d] was ignored: Empty job name.", i)
			continue
		}

		je := NewJobEx()
		je.Node = record[nodeIdx]
		if port, err := strconv.Atoi(record[portIdx]); err == nil {
			je.Port = port
		}
		je.FilePath = record[pathIdx]
		je.Param = record[paramIdx]
		je.Env = record[envIdx]
		je.Workspace = record[workIdx]
		if wrc, err := strconv.Atoi(record[wrcIdx]); err == nil {
			je.WrnRC = wrc
		}
		je.WrnPtn = record[wptIdx]
		if erc, err := strconv.Atoi(record[ercIdx]); err == nil {
			je.ErrRC = erc
		}
		je.ErrPtn = record[eptIdx]
		if tmout, err := strconv.Atoi(record[tmoutIdx]); err == nil {
			je.TimeoutMin = tmout
		}

		if len(record) >= withSecondary {
			je.SecondaryNode = record[secNodeIdx]
			if port, err := strconv.Atoi(record[secPortIdx]); err == nil {
				je.SecondaryPort = port
			}
		}

		jobExMap[name] = je
	}

	return jobExMap, nil
}
