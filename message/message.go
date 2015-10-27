// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package message

// 通信メッセージの共通インタフェース
type Message interface {
	ParseJSON(string) error
	GenerateJSON() (string, error)
}

const HEARTBEAT = "heartbeat"
const ST_HEADER = "ST:"
