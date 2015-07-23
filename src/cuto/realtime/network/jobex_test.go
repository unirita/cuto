package network

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestLoadJobexFromReader(t *testing.T) {
	csv := `
ジョブ名,ノード名,ポート番号,実行ファイル,パラメータ,環境変数,作業フォルダ,警告コード,警告出力,異常コード,異常出力,タイムアウト,セカンダリ実行ノード,セカンダリポート番号
job1,123.45.67.89,1234,/scripts/job1.sh,param1,env1,/work,10,warn1,20,err1,3600,secondary,12345
job2,12.345.67.89,5678,/scripts/job2.sh,param2,env2,/work2,11,warn2,21,err2,3600,,`

	jobex, err := LoadJobexFromReader(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if len(jobex) != 3 {
		t.Errorf("len(jobex) => %d, want %d", len(jobex), 3)
	}
}

func TestMergeParamIntoJobex(t *testing.T) {
	base := [][]string{
		[]string{"job1", "node1", "1234", "/scripts/job1.sh", "param1", "env1", "/work", "10", "warn1", "20", "err1", "3600", "secondary", "12345"},
		[]string{"job2", "node2", "2345", "/scripts/job2.sh", "param2", "env2", "/work2", "11", "warn2", "21", "err2", "3600", "secondary2", "23456"},
	}

	job1 := Job{}
	job1.Name = new(string)
	*job1.Name = "job1"
	job1.Node = new(string)
	*job1.Node = "realnode"
	job1.Port = new(int)
	*job1.Port = 1000
	job1.Path = new(string)
	*job1.Path = "/scripts/real.sh"
	job1.Param = new(string)
	*job1.Param = "realparam"
	job1.Env = new(string)
	*job1.Env = "realenv"
	job1.Work = new(string)
	*job1.Work = "/realwork"
	job1.WRC = new(int)
	*job1.WRC = 100
	job1.WPtn = new(string)
	*job1.WPtn = "realwrn"
	job1.ERC = new(int)
	*job1.ERC = 110
	job1.EPtn = new(string)
	*job1.EPtn = "realerr"
	job1.Timeout = new(int)
	*job1.Timeout = 1800
	job1.SNode = new(string)
	*job1.SNode = "realsnode"
	job1.SPort = new(int)
	*job1.SPort = 2000
	job2 := Job{}
	job2.Name = new(string)
	*job2.Name = "job2"
	job2.Node = new(string)
	*job2.Node = "realnode2"
	job3 := Job{}
	job3.Name = new(string)
	*job3.Name = "job3"
	job3.Node = new(string)
	*job3.Node = "realnode3"
	job3.Port = new(int)
	*job3.Port = 1002
	job3.Path = new(string)
	*job3.Path = "/scripts/real3.sh"
	job3.Param = new(string)
	*job3.Param = "realparam3"
	job3.Env = new(string)
	*job3.Env = "realenv3"
	job3.Work = new(string)
	*job3.Work = "/realwork3"
	job3.WRC = new(int)
	*job3.WRC = 102
	job3.WPtn = new(string)
	*job3.WPtn = "realwrn3"
	job3.ERC = new(int)
	*job3.ERC = 112
	job3.EPtn = new(string)
	*job3.EPtn = "realerr3"
	job3.Timeout = new(int)
	*job3.Timeout = 1802
	job3.SNode = new(string)
	*job3.SNode = "realsnode3"
	job3.SPort = new(int)
	*job3.SPort = 2002
	jobs := []Job{job1, job2, job3}

	baseJob2 := make([]string, 14)
	copy(baseJob2, base[1])
	result := MergeParamIntoJobex(base, jobs)
	if len(result) != 3 {
		t.Fatalf("len(result) => %d, want %d", len(result), 3)
	}
	if result[0][nodeIdx] != *job1.Node {
		t.Errorf("result[0][nodeIdx] => %s, want %s", result[0][nodeIdx], *job1.Node)
	}
	if result[0][portIdx] != strconv.Itoa(*job1.Port) {
		t.Errorf("result[0][portIdx] => %s, want %d", result[0][portIdx], *job1.Port)
	}
	if result[0][pathIdx] != *job1.Path {
		t.Errorf("result[0][pathIdx] => %s, want %s", result[0][pathIdx], *job1.Path)
	}
	if result[0][paramIdx] != *job1.Param {
		t.Errorf("result[0][paramIdx] => %s, want %s", result[0][paramIdx], *job1.Param)
	}
	if result[0][envIdx] != *job1.Env {
		t.Errorf("result[0][envIdx] => %s, want %s", result[0][envIdx], *job1.Env)
	}
	if result[0][workIdx] != *job1.Work {
		t.Errorf("result[0][workIdx] => %s, want %s", result[0][workIdx], *job1.Work)
	}
	if result[0][wrcIdx] != strconv.Itoa(*job1.WRC) {
		t.Errorf("result[0][wrcIdx] => %s, want %d", result[0][wrcIdx], *job1.WRC)
	}
	if result[0][wptnIdx] != *job1.WPtn {
		t.Errorf("result[0][wptnIdx] => %s, want %s", result[0][wptnIdx], *job1.WPtn)
	}
	if result[0][ercIdx] != strconv.Itoa(*job1.ERC) {
		t.Errorf("result[0][ercIdx] => %s, want %d", result[0][ercIdx], *job1.ERC)
	}
	if result[0][eptnIdx] != *job1.EPtn {
		t.Errorf("result[0][eptnIdx] => %s, want %s", result[0][eptnIdx], *job1.EPtn)
	}
	if result[0][timeoutIdx] != strconv.Itoa(*job1.Timeout) {
		t.Errorf("result[0][timeoutIdx] => %s, want %d", result[0][timeoutIdx], *job1.Timeout)
	}
	if result[0][snodeIdx] != *job1.SNode {
		t.Errorf("result[0][snodeIdx] => %s, want %s", result[0][snodeIdx], *job1.SNode)
	}
	if result[0][sportIdx] != strconv.Itoa(*job1.SPort) {
		t.Errorf("result[0][sportIdx] => %s, want %d", result[0][sportIdx], *job1.SPort)
	}
	if result[1][nodeIdx] != *job2.Node {
		t.Errorf("result[1][nodeIdx] => %s, want %s", result[0][nodeIdx], *job2.Node)
	}
	if !reflect.DeepEqual(result[1][2:], baseJob2[2:]) {
		t.Errorf("Unexpected properties changed in job2.")
	}
	if result[2][nameIdx] != *job3.Name {
		t.Errorf("result[2][nameIdx] => %s, want %s", result[2][nameIdx], *job3.Name)
	}
	if result[2][nodeIdx] != *job3.Node {
		t.Errorf("result[2][nodeIdx] => %s, want %s", result[2][nodeIdx], *job3.Node)
	}
	if result[2][portIdx] != strconv.Itoa(*job3.Port) {
		t.Errorf("result[2][portIdx] => %s, want %d", result[2][portIdx], *job3.Port)
	}
	if result[2][pathIdx] != *job3.Path {
		t.Errorf("result[2][pathIdx] => %s, want %s", result[2][pathIdx], *job3.Path)
	}
	if result[2][paramIdx] != *job3.Param {
		t.Errorf("result[2][paramIdx] => %s, want %s", result[2][paramIdx], *job3.Param)
	}
	if result[2][envIdx] != *job3.Env {
		t.Errorf("result[2][envIdx] => %s, want %s", result[2][envIdx], *job3.Env)
	}
	if result[2][workIdx] != *job3.Work {
		t.Errorf("result[2][workIdx] => %s, want %s", result[2][workIdx], *job3.Work)
	}
	if result[2][wrcIdx] != strconv.Itoa(*job3.WRC) {
		t.Errorf("result[2][wrcIdx] => %s, want %d", result[2][wrcIdx], *job3.WRC)
	}
	if result[2][wptnIdx] != *job3.WPtn {
		t.Errorf("result[2][wptnIdx] => %s, want %s", result[2][wptnIdx], *job3.WPtn)
	}
	if result[2][ercIdx] != strconv.Itoa(*job3.ERC) {
		t.Errorf("result[2][ercIdx] => %s, want %d", result[2][ercIdx], *job3.ERC)
	}
	if result[2][eptnIdx] != *job3.EPtn {
		t.Errorf("result[2][eptnIdx] => %s, want %s", result[2][eptnIdx], *job3.EPtn)
	}
	if result[2][timeoutIdx] != strconv.Itoa(*job3.Timeout) {
		t.Errorf("result[2][timeoutIdx] => %s, want %d", result[2][timeoutIdx], *job3.Timeout)
	}
	if result[2][snodeIdx] != *job3.SNode {
		t.Errorf("result[2][snodeIdx] => %s, want %s", result[2][snodeIdx], *job3.SNode)
	}
	if result[2][sportIdx] != strconv.Itoa(*job3.SPort) {
		t.Errorf("result[2][sportIdx] => %s, want %d", result[2][sportIdx], *job3.SPort)
	}
}
