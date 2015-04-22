// Copyright 2015 unirita Inc.
// Created 2015/04/20 T.Honda

package jobnet

import (
	"fmt"
	"testing"
)

type testJob struct {
	hasError   bool
	isExecuted bool
	Job
}

func (j *testJob) Execute() (Element, error) {
	j.isExecuted = true
	if j.hasError {
		return nil, fmt.Errorf("job[%s] has error.", j.id)
	}
	return j.Next, nil
}

func generateTestJob(idnum int) *testJob {
	j := new(testJob)
	j.id = fmt.Sprintf("jobid%d", idnum)
	j.Name = fmt.Sprintf("job%d", idnum)
	return j
}

func TestNewPath_分岐経路に先頭要素がセットされる(t *testing.T) {
	j1 := generateTestJob(1)

	p := NewPath(j1)
	if p.Head != j1 {
		t.Errorf("セットされた先頭要素[%v]は想定と違っている", p.Head)
	}
}

func TestPathRun_ゲートウェイまでのすべてのジョブが実行される_ゲートウェイの後続要素なし(t *testing.T) {
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	g1 := NewGateway("gwid1")

	j1.AddNext(j2)
	j2.AddNext(g1)
	p := NewPath(j1)

	done := make(chan struct{}, 1)
	p.Run(done)
	if p.Err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", p.Err)
	}
	if p.Goal != g1 {
		t.Fatalf("経路の終着ノード[%v]が想定と違っている。", p.Goal)
	}
	if !j1.isExecuted {
		t.Errorf("job1が実行されなかった。")
	}
	if !j2.isExecuted {
		t.Errorf("job2が実行されなかった。")
	}
}

func TestPathRun_ゲートウェイまでのすべてのジョブが実行される_ゲートウェイの後続要素あり(t *testing.T) {
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	g1 := NewGateway("gwid1")

	j1.AddNext(j2)
	j2.AddNext(g1)
	g1.AddNext(j3)
	p := NewPath(j1)

	done := make(chan struct{}, 1)
	p.Run(done)
	if p.Err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", p.Err)
	}
	if p.Goal != g1 {
		t.Fatalf("経路の終着ノード[%v]が想定と違っている。", p.Goal)
	}
	if !j1.isExecuted {
		t.Errorf("job1が実行されなかった。")
	}
	if !j2.isExecuted {
		t.Errorf("job2が実行されなかった。")
	}
	if j3.isExecuted {
		t.Errorf("経路外のジョブjob3が実行された。")
	}
}

func TestPathRun_ゲートウェイでない終端要素に到達したらエラー(t *testing.T) {
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)

	j1.AddNext(j2)
	p := NewPath(j1)

	done := make(chan struct{}, 1)
	p.Run(done)
	if p.Err == nil {
		t.Fatalf("エラーが発生していない。")
	}
	if p.Goal != nil {
		t.Fatalf("ゲートウェイに未到達にも関わらず、終着ノード[%v]がセットされてしまった。", p.Goal)
	}
	if !j1.isExecuted {
		t.Errorf("job1が実行されなかった。")
	}
	if !j2.isExecuted {
		t.Errorf("job2が実行されなかった。")
	}
}

func TestPathRun_経路の途中でジョブが異常終了したらエラー(t *testing.T) {
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	g1 := NewGateway("gwid1")

	j2.hasError = true
	j1.AddNext(j2)
	j2.AddNext(j3)
	j3.AddNext(g1)
	p := NewPath(j1)

	done := make(chan struct{}, 1)
	p.Run(done)
	if p.Err == nil {
		t.Fatalf("エラーが発生していない。")
	}
	if p.Goal != nil {
		t.Fatalf("ゲートウェイに未到達にも関わらず、終着ノード[%v]がセットされてしまった。", p.Goal)
	}
	if !j1.isExecuted {
		t.Errorf("job1が実行されなかった。")
	}
	if !j2.isExecuted {
		t.Errorf("job2が実行されなかった。")
	}
	if j3.isExecuted {
		t.Errorf("異常終了ジョブの後続ジョブjob3が実行された。")
	}
}
