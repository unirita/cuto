package network

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoadJobexFromReader(t *testing.T) {
	csv := `
ジョブ名,ノード名,ポート番号,実行ファイル,パラメータ,環境変数,作業フォルダ,警告コード,警告出力,異常コード,異常出力,タイムアウト,セカンダリ実行ノード,セカンダリポート番号
job1,123.45.67.89,1234,/scripts/job1.sh,param1,env1,/work,10,warn1,20,err1,3600,secondary,12345
job2,12.345.67.89,5678,/scripts/job2.sh,param2,env2,/work2,11,warn2,21,err2,3600,,`

	jobex, err := loadJobexFromReader(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if len(jobex) != 2 {
		t.Fatalf("len(jobex) => %d, want %d", len(jobex), 2)
	}
}

func TestExportJob(t *testing.T) {
	expected := `,,,,,,,,,,,,,
job1,node1,1234,/scripts/job1.sh,param1,env1,/work1,11,warn1,21,err1,100,snode1,1000
job2,node2,2345,/scripts/job2.sh,param2,env2,/work2,12,warn2,22,err2,200,snode2,2000
`

	n := new(Network)
	n.Jobs = []Job{
		Job{
			Name:    "job1",
			Node:    "node1",
			Port:    1234,
			Path:    "/scripts/job1.sh",
			Param:   "param1",
			Env:     "env1",
			Work:    "/work1",
			WRC:     11,
			WPtn:    "warn1",
			ERC:     21,
			EPtn:    "err1",
			Timeout: 100,
			SNode:   "snode1",
			SPort:   1000,
		},
		Job{
			Name:    "job2",
			Node:    "node2",
			Port:    2345,
			Path:    "/scripts/job2.sh",
			Param:   "param2",
			Env:     "env2",
			Work:    "/work2",
			WRC:     12,
			WPtn:    "warn2",
			ERC:     22,
			EPtn:    "err2",
			Timeout: 200,
			SNode:   "snode2",
			SPort:   2000,
		},
	}

	w := new(bytes.Buffer)
	err := n.exportJob(w)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	result := w.String()
	if result != expected {
		t.Log("Exported csv is not expected.")
		t.Log("Expected:")
		t.Log(expected)
		t.Log("Actual:")
		t.Log(result)
		t.Fail()
	}
}

func TestParse(t *testing.T) {
	jsonStr := `
{
	"flow":"job1->job2->[job3,job4->job5]->job6",
	"jobs":[
		{
			"name":"job1",
			"node":"testnode",
			"port":1234,
			"path":"/scripts/job1.sh",
			"param":"abc",
			"env":"env1=test",
			"work":"/work",
			"wrc":5,
			"wptn":"warning",
			"erc":10,
			"eptn":"error",
			"timeout":30,
			"snode":"secondary",
			"sport":2345
		}
	]
}
`
	network, err := Parse(jsonStr)
	if err != nil {
		t.Fatalf("Unexpected error occurd: %s", err)
	}
	if network.Flow != "job1->job2->[job3,job4->job5]->job6" {
		t.Logf("Flow => %s", network.Flow)
		t.Logf("Want %s", "job1->job2->[job3,job4->job5]->job6")
		t.Fail()
	}
	if len(network.Jobs) != 1 {
		t.Fatalf("len(Jobs) => %d, want %d", len(network.Jobs), 1)
	}

	job := network.Jobs[0]
	if job.Name != "job1" {
		t.Errorf("job.Name => %s, want %s", job.Name, "job1")
	}
	if job.Node != "testnode" {
		t.Errorf("job.Node => %s, want %s", job.Node, "testnode")
	}
	if job.Port != 1234 {
		t.Errorf("job.Port => %d, want %d", job.Port, 1234)
	}
	if job.Path != "/scripts/job1.sh" {
		t.Errorf("job.Path => %s, want %s", job.Path, "/scripts/job1.sh")
	}
	if job.Param != "abc" {
		t.Errorf("job.Param => %s, want %s", job.Param, "abc")
	}
	if job.Env != "env1=test" {
		t.Errorf("job.Env => %s, want %s", job.Env, "env1=test")
	}
	if job.Work != "/work" {
		t.Errorf("job.Work => %s, want %s", job.Work, "/work")
	}
	if job.WRC != 5 {
		t.Errorf("job.WRC => %d, want %d", job.WRC, 5)
	}
	if job.WPtn != "warning" {
		t.Errorf("job.WPtn => %s, want %s", job.WPtn, "warning")
	}
	if job.ERC != 10 {
		t.Errorf("job.ERC => %d, want %d", job.ERC, 10)
	}
	if job.EPtn != "error" {
		t.Errorf("job.EPtn => %s, want %s", job.EPtn, "error")
	}
	if job.Timeout != 30 {
		t.Errorf("job.Timeout => %d, want %d", job.Timeout, 30)
	}
	if job.SNode != "secondary" {
		t.Errorf("job.SNode => %s, want %s", job.SNode, "secondary")
	}
	if job.SPort != 2345 {
		t.Errorf("job.SPort => %d, want %d", job.SPort, 2345)
	}
}

func TestParse_WithNullValue(t *testing.T) {
	jsonStr := `
{
	"flow":"job1->job2->[job3,job4->job5]->job6",
	"jobs":[
		{
			"name":"job1"
		},
		{
			"name":"job2"
		},
		{
			"name":"job3",
			"node":"realtimenode"
		}
	]
}
`
	jobex = [][]string{
		[]string{
			"job2",
			"node2",
			"2345",
			"/scripts/job2.sh",
			"param2",
			"env2",
			"/work2",
			"12",
			"warn2",
			"22",
			"err2",
			"200",
			"snode2",
			"2000",
		},
		[]string{
			"job3",
			"node3",
			"3456",
			"/scripts/job3.sh",
			"param3",
			"env3",
			"/work3",
			"13",
			"warn3",
			"23",
			"err3",
			"300",
			"snode3",
			"3000",
		},
	}
	defer func() {
		jobex := make([][]string, 1)
		jobex[0] = make([]string, columns)
	}()

	network, err := Parse(jsonStr)
	if err != nil {
		t.Fatalf("Unexpected error occurd: %s", err)
	}
	if len(network.Jobs) != 3 {
		t.Fatalf("len(Jobs) => %d, want %d", len(network.Jobs), 3)
	}

	firstJob := network.Jobs[0]
	if firstJob.Name != "job1" {
		t.Errorf("firstJob.Name => %s, want %s", firstJob.Name, "job1")
	}
	if firstJob.Node != "" {
		t.Errorf("firstJob.Node => %s, want %s", firstJob.Node, "")
	}
	if firstJob.Port != 0 {
		t.Errorf("firstJob.Port => %d, want %d", firstJob.Port, 0)
	}
	if firstJob.Path != "" {
		t.Errorf("firstJob.Path => %s, want %s", firstJob.Path, "")
	}
	if firstJob.Param != "" {
		t.Errorf("firstJob.Param => %s, want %s", firstJob.Param, "")
	}
	if firstJob.Env != "" {
		t.Errorf("firstJob.Env => %s, want %s", firstJob.Env, "")
	}
	if firstJob.Work != "" {
		t.Errorf("firstJob.Work => %s, want %s", firstJob.Work, "")
	}
	if firstJob.WRC != 0 {
		t.Errorf("firstJob.WRC => %d, want %d", firstJob.WRC, 0)
	}
	if firstJob.WPtn != "" {
		t.Errorf("firstJob.WPtn => %s, want %s", firstJob.WPtn, "")
	}
	if firstJob.ERC != 0 {
		t.Errorf("firstJob.ERC => %d, want %d", firstJob.ERC, 0)
	}
	if firstJob.EPtn != "" {
		t.Errorf("firstJob.EPtn => %s, want %s", firstJob.EPtn, "")
	}
	if firstJob.Timeout != 0 {
		t.Errorf("firstJob.Timeout => %d, want %d", firstJob.Timeout, 0)
	}
	if firstJob.SNode != "" {
		t.Errorf("firstJob.SNode => %s, want nil", firstJob.SNode)
	}
	if firstJob.SPort != 0 {
		t.Errorf("firstJob.SPort => %d, want nil", firstJob.SPort)
	}

	secondJob := network.Jobs[1]
	if secondJob.Name != "job2" {
		t.Errorf("secondJob.Name => %s, want %s", secondJob.Name, "job2")
	}
	if secondJob.Node != "node2" {
		t.Errorf("secondJob.Node => %s, want %s", secondJob.Node, "node2")
	}
	if secondJob.Port != 2345 {
		t.Errorf("secondJob.Port => %d, want %d", secondJob.Port, 2345)
	}
	if secondJob.Path != "/scripts/job2.sh" {
		t.Errorf("secondJob.Path => %s, want %s", secondJob.Path, "/scripts/job2.sh")
	}
	if secondJob.Param != "param2" {
		t.Errorf("secondJob.Param => %s, want %s", secondJob.Param, "param2")
	}
	if secondJob.Env != "env2" {
		t.Errorf("secondJob.Env => %s, want %s", secondJob.Env, "env2")
	}
	if secondJob.Work != "/work2" {
		t.Errorf("secondJob.Work => %s, want %s", secondJob.Work, "/work2")
	}
	if secondJob.WRC != 12 {
		t.Errorf("secondJob.WRC => %d, want %d", secondJob.WRC, 12)
	}
	if secondJob.WPtn != "warn2" {
		t.Errorf("secondJob.WPtn => %s, want %s", secondJob.WPtn, "warn2")
	}
	if secondJob.ERC != 22 {
		t.Errorf("secondJob.ERC => %d, want %d", secondJob.ERC, 22)
	}
	if secondJob.EPtn != "err2" {
		t.Errorf("secondJob.EPtn => %s, want %s", secondJob.EPtn, "err2")
	}
	if secondJob.Timeout != 200 {
		t.Errorf("secondJob.Timeout => %d, want %d", secondJob.Timeout, 200)
	}
	if secondJob.SNode != "snode2" {
		t.Errorf("secondJob.SNode => %s, want %s", secondJob.SNode, "snode2")
	}
	if secondJob.SPort != 2000 {
		t.Errorf("secondJob.SPort => %d, want %d", secondJob.SPort, 2000)
	}

	thirdJob := network.Jobs[2]
	if thirdJob.Name != "job3" {
		t.Errorf("thirdJob.Name => %s, want %s", thirdJob.Name, "job3")
	}
	if thirdJob.Node != "realtimenode" {
		t.Errorf("thirdJob.Node => %s, want %s", thirdJob.Node, "realtimenode")
	}
	if thirdJob.Port != 3456 {
		t.Errorf("thirdJob.Port => %d, want %d", thirdJob.Port, 3456)
	}
	if thirdJob.Path != "/scripts/job3.sh" {
		t.Errorf("thirdJob.Path => %s, want %s", thirdJob.Path, "/scripts/job3.sh")
	}
	if thirdJob.Param != "param3" {
		t.Errorf("thirdJob.Param => %s, want %s", thirdJob.Param, "param3")
	}
	if thirdJob.Env != "env3" {
		t.Errorf("thirdJob.Env => %s, want %s", thirdJob.Env, "env3")
	}
	if thirdJob.Work != "/work3" {
		t.Errorf("thirdJob.Work => %s, want %s", thirdJob.Work, "/work3")
	}
	if thirdJob.WRC != 13 {
		t.Errorf("thirdJob.WRC => %d, want %d", thirdJob.WRC, 13)
	}
	if thirdJob.WPtn != "warn3" {
		t.Errorf("thirdJob.WPtn => %s, want %s", thirdJob.WPtn, "warn3")
	}
	if thirdJob.ERC != 23 {
		t.Errorf("thirdJob.ERC => %d, want %d", thirdJob.ERC, 23)
	}
	if thirdJob.EPtn != "err3" {
		t.Errorf("thirdJob.EPtn => %s, want %s", thirdJob.EPtn, "err3")
	}
	if thirdJob.Timeout != 300 {
		t.Errorf("thirdJob.Timeout => %d, want %d", thirdJob.Timeout, 300)
	}
	if thirdJob.SNode != "snode3" {
		t.Errorf("thirdJob.SNode => %s, want %s", thirdJob.SNode, "snode3")
	}
	if thirdJob.SPort != 3000 {
		t.Errorf("thirdJob.SPort => %d, want %d", thirdJob.SPort, 3000)
	}
}

func TestParse_NonParameterJobex(t *testing.T) {
	jsonStr := `
{
	"flow":"job1->job2->[job3,job4->job5]->job6",
	"jobs":[]
}
`
	jobex = [][]string{
		[]string{
			"job2",
			"node2",
			"2345",
			"/scripts/job2.sh",
			"param2",
			"env2",
			"/work2",
			"12",
			"warn2",
			"22",
			"err2",
			"200",
			"snode2",
			"2000",
		},
	}
	defer func() {
		jobex := make([][]string, 1)
		jobex[0] = make([]string, columns)
	}()

	network, err := Parse(jsonStr)
	if err != nil {
		t.Fatalf("Unexpected error occurd: %s", err)
	}
	if len(network.Jobs) != 1 {
		t.Fatalf("len(Jobs) => %d, want %d", len(network.Jobs), 1)
	}

	job := network.Jobs[0]
	if job.Name != "job2" {
		t.Errorf("job.Name => %s, want %s", job.Name, "job2")
	}
	if job.Node != "node2" {
		t.Errorf("job.Node => %s, want %s", job.Node, "node2")
	}
	if job.Port != 2345 {
		t.Errorf("job.Port => %d, want %d", job.Port, 2345)
	}
	if job.Path != "/scripts/job2.sh" {
		t.Errorf("job.Path => %s, want %s", job.Path, "/scripts/job2.sh")
	}
	if job.Param != "param2" {
		t.Errorf("job.Param => %s, want %s", job.Param, "param2")
	}
	if job.Env != "env2" {
		t.Errorf("job.Env => %s, want %s", job.Env, "env2")
	}
	if job.Work != "/work2" {
		t.Errorf("job.Work => %s, want %s", job.Work, "/work2")
	}
	if job.WRC != 12 {
		t.Errorf("job.WRC => %d, want %d", job.WRC, 12)
	}
	if job.WPtn != "warn2" {
		t.Errorf("job.WPtn => %s, want %s", job.WPtn, "warn2")
	}
	if job.ERC != 22 {
		t.Errorf("job.ERC => %d, want %d", job.ERC, 22)
	}
	if job.EPtn != "err2" {
		t.Errorf("job.EPtn => %s, want %s", job.EPtn, "err2")
	}
	if job.Timeout != 200 {
		t.Errorf("job.Timeout => %d, want %d", job.Timeout, 200)
	}
	if job.SNode != "snode2" {
		t.Errorf("job.SNode => %s, want %s", job.SNode, "snode2")
	}
	if job.SPort != 2000 {
		t.Errorf("job.SPort => %d, want %d", job.SPort, 2000)
	}
}

func TestParse_WithJSONError(t *testing.T) {
	jsonStr := `
{
	flow:"job1->job2->[job3,job4->job5]->job6",
	"jobs":[
		{
			"name":"job2",
			"node":"testnode",
			"port":1234,
			"path":"/scripts/job2.sh",
			"param":"abc",
			"env":"env1=test",
			"work":"/work",
			"wrc":5,
			"wptn":"warning",
			"erc":10,
			"eptn":"error",
			"timeout":30,
			"snode":"secondary",
			"sport":2345
		},
		{
			"name":"job5",
			"param":"bcd"
		}
	]
}
`
	_, err := Parse(jsonStr)
	if err == nil {
		t.Fatalf("No error occured.")
	}
}

func TestParse_WithAnonymousJob(t *testing.T) {
	jsonStr := `
{
	"flow":"job1->job2->[job3,job4->job5]->job6",
	"jobs":[
		{
			"node":"testnode",
			"port":1234,
			"path":"/scripts/job2.sh",
			"param":"abc",
			"env":"env1=test",
			"work":"/work",
			"wrc":5,
			"wptn":"warning",
			"erc":10,
			"eptn":"error",
			"timeout":30,
			"snode":"secondary",
			"sport":2345
		},
		{
			"name":"job5",
			"param":"bcd"
		}
	]
}
`
	_, err := Parse(jsonStr)
	if err == nil {
		t.Fatalf("No error occured.")
	}
}
