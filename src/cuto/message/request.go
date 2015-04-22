// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package message

import (
	"encoding/json"
	"fmt"
)

// 要求メッセージ。
type Request struct {
	Type      string `json:"type"`
	NID       int    `json:"nid"`
	JID       string `json:"jid"`
	Path      string `json:"path"`
	Param     string `json:"param"`
	Env       string `json:"env"`
	Workspace string `json:"workspace"`
	WarnRC    int    `json:"warnrc"`
	WarnStr   string `json:"warnstr"`
	ErrRC     int    `json:"errrc"`
	ErrStr    string `json:"errstr"`
	Timeout   int    `json:"timeout"`
}

const requestMessageType = "request"

// ジョブ実行要求JSONメッセージをパースし、Requestオブジェクトのメンバをセットする。
//
// param : message 受信メッセージ文字列。
func (r *Request) ParseJSON(message string) error {
	byteMessage := []byte(message)

	err := json.Unmarshal(byteMessage, r)

	if err != nil {
		return err
	}

	if r.Type != requestMessageType {
		return fmt.Errorf("Invalid message type.")
	}

	return nil
}

// Requestオブジェクトの値を元に、ジョブ実行要求JSONメッセージを生成する
//
// return : JSONメッセージフォーマットの文字列。
func (r Request) GenerateJSON() (string, error) {
	r.Type = requestMessageType
	byteMessage, err := json.Marshal(r)
	if err != nil {
		return ``, err
	}

	return string(byteMessage), nil
}

// masterで利用可能な変数を展開する。
func (r *Request) ExpandMasterVars() error {
	newPath, err := ExpandStringVars(r.Path, plcMaster, kndEnv)
	if err != nil {
		return err
	}
	newParam, err := ExpandStringVars(r.Param, plcMaster, kndSys, kndEnv, kndJob)
	if err != nil {
		return err
	}
	newEnv, err := ExpandStringVars(r.Env, plcMaster, kndSys, kndEnv, kndJob)
	if err != nil {
		return err
	}
	newWork, err := ExpandStringVars(r.Workspace, plcMaster, kndEnv)
	if err != nil {
		return err
	}

	// 全て展開成功した場合のみメッセージの内容を書き換えする
	r.Path = newPath
	r.Param = newParam
	r.Env = newEnv
	r.Workspace = newWork
	return nil
}

// servantで利用可能な変数を展開する。
func (r *Request) ExpandServantVars() error {
	newPath, err := ExpandStringVars(r.Path, plcServant, kndSys, kndEnv)
	if err != nil {
		return err
	}
	newParam, err := ExpandStringVars(r.Param, plcServant, kndSys, kndEnv)
	if err != nil {
		return err
	}
	newEnv, err := ExpandStringVars(r.Env, plcServant, kndSys, kndEnv)
	if err != nil {
		return err
	}
	newWork, err := ExpandStringVars(r.Workspace, plcServant, kndSys, kndEnv)
	if err != nil {
		return err
	}

	// 全て展開成功した場合のみメッセージの内容を書き換えする
	r.Path = newPath
	r.Param = newParam
	r.Env = newEnv
	r.Workspace = newWork
	return nil
}
