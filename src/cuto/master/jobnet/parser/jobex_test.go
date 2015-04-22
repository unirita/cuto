package parser

import (
	"strings"
	"testing"
)

func TestParseJobExFile_ファイルが存在しない場合は空のマップを返す(t *testing.T) {
	jeMap, err := ParseJobExFile("noexistsfilepath")
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if jeMap == nil {
		t.Fatal("空のマップではなく、nilが返された。")
	}
	if len(jeMap) != 0 {
		t.Error("返されたマップが空ではない。")
	}
}

func TestParseJobEx_拡張ジョブ定義CSVをパースできる(t *testing.T) {
	csv := `
ジョブ名,ノード名,ポート番号,実行ファイル,パラメータ,環境変数,作業フォルダ,警告コード,警告出力,異常コード,異常出力,タイムアウト
testjob1,123.45.67.89,1234,C:\work\test1.bat,testparam1,testenv1,C:\work1,10,warn1,11,err1,3600
testjob2,12.345.67.89,5678,C:\work\test2.bat,testparam2,testenv2,C:\work2,20,warn2,21,err2,3600`

	r := strings.NewReader(csv)
	jeMap, err := ParseJobEx(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	if len(jeMap) != 2 {
		t.Fatalf("パース結果が2件になるはずが、%d件になった。", len(jeMap))
	}

	j1, ok := jeMap["testjob1"]
	if !ok {
		t.Fatalf("パース結果にtestjob1がセットされていない。")
	}

	if j1.Node != `123.45.67.89` {
		t.Errorf("testjob1のノード名のパース結果[%s]が間違っています。", j1.Node)
	}
	if j1.Port != 1234 {
		t.Errorf("testjob1のポート番号のパース結果[%d]が間違っています。", j1.Port)
	}
	if j1.FilePath != `C:\work\test1.bat` {
		t.Errorf("testjob1の実行ファイルパスのパース結果[%s]が間違っています。", j1.FilePath)
	}
	if j1.Param != `testparam1` {
		t.Errorf("testjob1の実行時パラメータのパース結果[%s]が間違っています。", j1.Param)
	}
	if j1.Env != `testenv1` {
		t.Errorf("testjob1の環境変数のパース結果[%s]が間違っています。", j1.Env)
	}
	if j1.WrnRC != 10 {
		t.Errorf("testjob1の警告条件コードのパース結果[%d]が間違っています。", j1.WrnRC)
	}
	if j1.WrnPtn != `warn1` {
		t.Errorf("testjob1の警告出力パターンのパース結果[%s]が間違っています。", j1.WrnPtn)
	}
	if j1.ErrRC != 11 {
		t.Errorf("testjob1の異常条件コードのパース結果[%d]が間違っています。", j1.ErrRC)
	}
	if j1.ErrPtn != `err1` {
		t.Errorf("testjob1の異常出力パターンのパース結果[%s]が間違っています。", j1.ErrPtn)
	}
	if j1.TimeoutMin != 3600 {
		t.Errorf("testjob1の実行タイムアウト時間のパース結果[%d]が間違っています。", j1.TimeoutMin)
	}

	if _, ok := jeMap["testjob2"]; !ok {
		t.Fatalf("パース結果にtestjob2がセットされていない。")
	}
}

func TestParseJobEx_空のカラムにゼロ値がセットされる(t *testing.T) {
	csv := `
ジョブ名,ノード名,ポート番号,実行ファイル,パラメータ,環境変数,作業フォルダ,警告コード,警告出力,異常コード,異常出力,タイムアウト
testjob1,,,,,,,,,,,`

	r := strings.NewReader(csv)
	jeMap, err := ParseJobEx(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	if len(jeMap) != 1 {
		t.Fatalf("パース結果が1件になるはずが、%d件になった。", len(jeMap))
	}

	j1, ok := jeMap["testjob1"]
	if !ok {
		t.Fatalf("パース結果にtestjob1がセットされていない。")
	}

	if j1.Node != `` {
		t.Errorf("testjob1のノード名のパース結果[%s]が間違っています。", j1.Node)
	}
	if j1.Port != 0 {
		t.Errorf("testjob1のポート番号のパース結果[%d]が間違っています。", j1.Port)
	}
	if j1.FilePath != `` {
		t.Errorf("testjob1の実行ファイルパスのパース結果[%s]が間違っています。", j1.FilePath)
	}
	if j1.Param != `` {
		t.Errorf("testjob1の実行時パラメータのパース結果[%s]が間違っています。", j1.Param)
	}
	if j1.Env != `` {
		t.Errorf("testjob1の環境変数のパース結果[%s]が間違っています。", j1.Env)
	}
	if j1.WrnRC != 0 {
		t.Errorf("testjob1の警告条件コードのパース結果[%s]が間違っています。", j1.WrnRC)
	}
	if j1.WrnPtn != `` {
		t.Errorf("testjob1の警告出力パターンのパース結果[%s]が間違っています。", j1.WrnPtn)
	}
	if j1.ErrRC != 0 {
		t.Errorf("testjob1の異常条件コードのパース結果[%d]が間違っています。", j1.ErrRC)
	}
	if j1.ErrPtn != `` {
		t.Errorf("testjob1の異常出力パターンのパース結果[%s]が間違っています。", j1.ErrPtn)
	}
	if j1.TimeoutMin != -1 {
		t.Errorf("testjob1の実行タイムアウト時間のパース結果[%d]が間違っています。", j1.TimeoutMin)
	}
}

func TestParseJobEx_ジョブ名が無い行はエラーとして無視する(t *testing.T) {
	csv := `
ジョブ名,ノード名,ポート番号,実行ファイル,パラメータ,環境変数,作業フォルダ,警告コード,警告出力,異常コード,異常出力,タイムアウト
,123.45.67.89,1234,C:\work\test1.bat,testparam1,testenv1,C:\work1,10,warn1,11,err1,3600
testjob2,12.345.67.89,5678,C:\work\test2.bat,testparam2,testenv2,C:\work2,20,warn2,21,err2,3600`

	r := strings.NewReader(csv)
	jeMap, err := ParseJobEx(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	if len(jeMap) != 1 {
		t.Fatalf("パース結果が2件になるはずが、%d件になった。", len(jeMap))
	}

	if _, ok := jeMap["testjob2"]; !ok {
		t.Fatalf("パース結果にtestjob2がセットされていない。")
	}
}
