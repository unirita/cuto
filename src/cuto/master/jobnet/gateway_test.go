// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package jobnet

import "testing"

// ※補足事項
// testJob構造体およびその生成関数generateTestJobはpath_test.goで定義されています。

func TestNewGateway_ゲートウェイのIDがセットされる(t *testing.T) {
	g := NewGateway("gwid1")
	if g.ID() != "gwid1" {
		t.Errorf("セットされたID[%s]が想定と違っている。", g.ID())
	}
}

func TestGatewayType_ゲートウェイのノードタイプを取得できる(t *testing.T) {
	g := NewGateway("gwid1")
	if g.Type() != ELM_GW {
		t.Errorf("取得したノードタイプ[%d]が想定と違っている。", g.Type())
	}
}

func TestGatewayAddNext_ゲートウェイ構造体に後続エレメントを追加できる(t *testing.T) {
	g1 := new(Gateway)
	g2 := new(Gateway)
	j1 := new(Job)

	g1.AddNext(j1)
	if len(g1.Nexts) != 1 {
		t.Fatal("後続ジョブの追加に失敗した。")
	}
	if g1.Nexts[0] != j1 {
		t.Error("追加した後続ジョブが間違っている。")
	}

	g1.AddNext(g2)
	if len(g1.Nexts) != 2 {
		t.Fatal("後続ゲートウェイの追加に失敗した。")
	}
	if g1.Nexts[0] != j1 {
		t.Error("無関係な後続エレメントを変更している。")
	}
	if g1.Nexts[1] != g2 {
		t.Error("追加した後続ゲートウェイが間違っている。")
	}
}

func TestGatewayHasNext_ゲートウェイ構造体の後続エレメントの有無をチェックできる(t *testing.T) {
	g1 := new(Gateway)
	g2 := new(Gateway)
	j1 := new(Job)

	if g1.HasNext() {
		t.Error("後続エレメントが無いのにも関わらず、HasNextがtrueを返した")
	}

	g1.Nexts = append(g1.Nexts, g1)
	if !g1.HasNext() {
		t.Error("後続ゲートウェイがあるのにも関わらず、HasNextがfalseを返した")
	}

	g2.Nexts = append(g2.Nexts, j1)
	if !g2.HasNext() {
		t.Error("後続ジョブがあるのにも関わらず、HasNextがfalseを返した")
	}

	g1.Nexts = append(g1.Nexts, j1)
	if !g1.HasNext() {
		t.Error("後続エレメントが複数あるのにも関わらず、HasNextがfalseを返した")
	}
}

func TestGatewayExecute_後続ノードが存在しないケース(t *testing.T) {
	g := NewGateway("gwid1")

	next, err := g.Execute()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if next != nil {
		t.Errorf("存在しないはずの後続ノード[%s]が取得された。", next.ID())
	}
}

func TestGatewayExecute_後続ノードが1つのケース(t *testing.T) {
	g := NewGateway("gwid1")
	j := generateTestJob(1)
	g.AddNext(j)

	next, err := g.Execute()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if next != j {
		t.Errorf("想定外の後続ノード[%s]が取得された。", next.ID())
	}
}

func TestGatewayExecute_後続ノードが複数のケース(t *testing.T) {
	g1 := NewGateway("gwid1")
	g2 := NewGateway("gwid2")
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	g1.AddNext(j1)
	g1.AddNext(j3)
	j1.AddNext(j2)
	j2.AddNext(g2)
	j3.AddNext(g2)

	next, err := g1.Execute()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if next != g2 {
		t.Errorf("想定外の後続ノード[%s]が取得された。", next.ID())
	}
	if !j1.isExecuted {
		t.Errorf("分岐経路中のジョブj1が実行されなかった。")
	}
	if !j2.isExecuted {
		t.Errorf("分岐経路中のジョブj2が実行されなかった。")
	}
	if !j3.isExecuted {
		t.Errorf("分岐経路中のジョブj3が実行されなかった。")
	}
}

func TestGatewayExecute_分岐経路中でジョブが異常終了したらエラー(t *testing.T) {
	g1 := NewGateway("gwid1")
	g2 := NewGateway("gwid2")
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	g1.AddNext(j1)
	g1.AddNext(j3)
	j1.AddNext(j2)
	j2.AddNext(g2)
	j3.AddNext(g2)

	j1.hasError = true
	next, err := g1.Execute()
	if err == nil {
		t.Fatalf("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("想定外の後続ノード[%s]が取得された。", next.ID())
	}
	if !j1.isExecuted {
		t.Errorf("分岐経路中のジョブj1が実行されなかった。")
	}
	if j2.isExecuted {
		t.Errorf("異常終了ジョブの後続ジョブj2が実行された。")
	}
	if !j3.isExecuted {
		t.Errorf("分岐経路中のジョブj3が実行されなかった。")
	}
}

func TestGatewayExecute_経路間で終着ノードが異なる場合はエラー(t *testing.T) {
	g1 := NewGateway("gwid1")
	g2 := NewGateway("gwid2")
	g3 := NewGateway("gwid3")
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	g1.AddNext(j1)
	g1.AddNext(j2)
	g1.AddNext(j3)
	j1.AddNext(g2)
	j2.AddNext(g2)
	j3.AddNext(g3)

	next, err := g1.Execute()
	if err == nil {
		t.Fatalf("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("想定外の後続ノード[%s]が取得された。", next.ID())
	}
}
