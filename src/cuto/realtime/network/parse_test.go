package network

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	jsonStr := `
{
	"flow":"job1->job2->[job3,job4->job5]->job6",
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
	network, err := Parse(strings.NewReader(jsonStr))
	if err != nil {
		t.Fatalf("Unexpected error occurd: %s", err)
	}
	if network.Flow != "job1->job2->[job3,job4->job5]->job6" {
		t.Logf("Flow => %s", network.Flow)
		t.Logf("Want %s", "job1->job2->[job3,job4->job5]->job6")
		t.Fail()
	}
	if len(network.Jobs) != 2 {
		t.Errorf("len(Jobs) => %d, want %d", len(network.Jobs), 2)
	}

	firstJob := network.Jobs[0]
	if *firstJob.Name != "job2" {
		t.Errorf("firstJob.Name => %s, want %s", *firstJob.Name, "job2")
	}
	if *firstJob.Node != "testnode" {
		t.Errorf("firstJob.Node => %s, want %s", *firstJob.Node, "testnode")
	}
	if *firstJob.Port != 1234 {
		t.Errorf("firstJob.Port => %d, want %d", *firstJob.Port, 1234)
	}
	if *firstJob.Path != "/scripts/job2.sh" {
		t.Errorf("firstJob.Path => %s, want %s", *firstJob.Path, "/scripts/job2.sh")
	}
	if *firstJob.Param != "abc" {
		t.Errorf("firstJob.Param => %s, want %s", *firstJob.Param, "abc")
	}
	if *firstJob.Env != "env1=test" {
		t.Errorf("firstJob.Env => %s, want %s", *firstJob.Env, "env1=test")
	}
	if *firstJob.Work != "/work" {
		t.Errorf("firstJob.Work => %s, want %s", *firstJob.Work, "/work")
	}
	if *firstJob.WRC != 5 {
		t.Errorf("firstJob.WRC => %d, want %d", *firstJob.WRC, 5)
	}
	if *firstJob.WPtn != "warning" {
		t.Errorf("firstJob.WPtn => %s, want %s", *firstJob.WPtn, "warning")
	}
	if *firstJob.ERC != 10 {
		t.Errorf("firstJob.ERC => %d, want %d", *firstJob.ERC, 10)
	}
	if *firstJob.EPtn != "error" {
		t.Errorf("firstJob.EPtn => %s, want %s", *firstJob.EPtn, "error")
	}
	if *firstJob.Timeout != 30 {
		t.Errorf("firstJob.Timeout => %d, want %d", *firstJob.Timeout, 30)
	}
	if *firstJob.SNode != "secondary" {
		t.Errorf("firstJob.SNode => %s, want %s", *firstJob.SNode, "secondary")
	}
	if *firstJob.SPort != 2345 {
		t.Errorf("firstJob.SPort => %d, want %d", *firstJob.SPort, 2345)
	}

	secondJob := network.Jobs[1]
	if *secondJob.Name != "job5" {
		t.Errorf("secondJob.Name => %s, want %s", *secondJob.Name, "job5")
	}
	if secondJob.Node != nil {
		t.Errorf("secondJob.Node => %s, want %s", *secondJob.Node, "")
	}
	if secondJob.Port != nil {
		t.Errorf("secondJob.Port => %d, want %d", *secondJob.Port, 0)
	}
	if secondJob.Path != nil {
		t.Errorf("secondJob.Path => %s, want %s", *secondJob.Path, "")
	}
	if *secondJob.Param != "bcd" {
		t.Errorf("secondJob.Param => %s, want %s", *secondJob.Param, "bcd")
	}
	if secondJob.Env != nil {
		t.Errorf("secondJob.Env => %s, want %s", *secondJob.Env, "")
	}
	if secondJob.Work != nil {
		t.Errorf("secondJob.Work => %s, want %s", *secondJob.Work, "")
	}
	if secondJob.WRC != nil {
		t.Errorf("secondJob.WRC => %d, want %d", *secondJob.WRC, 0)
	}
	if secondJob.WPtn != nil {
		t.Errorf("secondJob.WPtn => %s, want %s", *secondJob.WPtn, "")
	}
	if secondJob.ERC != nil {
		t.Errorf("secondJob.ERC => %d, want %d", *secondJob.ERC, 0)
	}
	if secondJob.EPtn != nil {
		t.Errorf("secondJob.EPtn => %s, want %s", *secondJob.EPtn, "")
	}
	if secondJob.Timeout != nil {
		t.Errorf("secondJob.Timeout => %d, want %d", *secondJob.Timeout, 0)
	}
	if secondJob.SNode != nil {
		t.Errorf("secondJob.SNode => %s, want nil", *secondJob.SNode)
	}
	if secondJob.SPort != nil {
		t.Errorf("secondJob.SPort => %d, want nil", *secondJob.SPort)
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
	_, err := Parse(strings.NewReader(jsonStr))
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
	_, err := Parse(strings.NewReader(jsonStr))
	if err == nil {
		t.Fatalf("No error occured.")
	}
}
