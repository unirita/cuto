// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package jobnet

// BPMNタグの要素インタフェース
type Element interface {
	ID() string                // ノードのID
	Type() elementType         // ノードのタイプ
	AddNext(e Element) error   // 次（後続）ノードを追加する。
	HasNext() bool             // 後続ノードの保持フラグ
	Execute() (Element, error) // ノードの処理を実行する。
}

type elementType int

const (
	ELM_JOB elementType = iota
	ELM_GW
)
