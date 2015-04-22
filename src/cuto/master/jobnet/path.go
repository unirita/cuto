// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package jobnet

import "fmt"

// 並列ジョブ実行の1経路を表す構造体。
type Path struct {
	Head Element
	Goal Element
	Err  error
}

// 新しい並列実行経路を生成する
func NewPath(head Element) *Path {
	p := new(Path)
	p.Head = head
	return p
}

// 並列実行経路を実行する。
// 実行完了後、doneチャンネルに空のstruct{}リテラルを入力する。
//
// param : done 通知用チャネル。
func (p *Path) Run(done chan<- struct{}) {
	current := p.Head

PATHLOOP:
	for {
		if current == nil {
			p.Err = fmt.Errorf("Network is terminated in branch.")
			break PATHLOOP
		}

		if current.Type() == ELM_GW {
			p.Goal = current
			break PATHLOOP
		}

		next, err := current.Execute()
		if err != nil {
			p.Err = err
			break PATHLOOP
		}

		current = next
	}

	done <- struct{}{}
}
