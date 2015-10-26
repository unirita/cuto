// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package message

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"cuto/utctime"
)

const (
	plcMaster  = 'M'
	plcServant = 'S'
	kndSys     = 'S'
	kndEnv     = 'E'
	kndJob     = 'J'
	kndTime    = 'T'
)

const minKeyLength = len(`$MEx$`)
const tagSeparator = `:`

// 変数名の解析結果を格納する構造体
type variable struct {
	key   string
	Place byte
	Kind  byte
	Name  string
	Tag   string
}

// ジョブネットワーク変数の値を格納する構造体
type jobValue struct {
	ID  string
	RC  string
	SD  string
	ED  string
	OUT string
}

var sysValues map[string]string
var jobValues map[string]*jobValue

func init() {
	sysValues = make(map[string]string)
	jobValues = make(map[string]*jobValue)
}

// システム変数の値を追加する。
func AddSysValue(name, tag, value string) {
	fullName := fmt.Sprintf("%s%s%s", name, tagSeparator, tag)
	sysValues[fullName] = value
}

// ジョブネットワーク変数の値を追加する。
func AddJobValue(name string, res *Response) {
	j := new(jobValue)
	j.ID = res.JID
	j.RC = strconv.Itoa(res.RC)
	j.SD = res.St
	j.ED = res.Et
	j.OUT = res.Var

	jobValues[name] = j
}

// 文字列src内の変数を展開する。
// 展開処理のパラメータとして場所識別子placeと利用可能種別kindsを指定する。
func ExpandStringVars(src string, place byte, kinds ...byte) (string, error) {
	if len(kinds) == 0 {
		return ``, fmt.Errorf("Invalid kind of variable.")
	}

	kindptn := string(kinds)
	pattern := fmt.Sprintf(`\$%c[%s](.+?)\$`, place, kindptn)
	exp := regexp.MustCompile(pattern)
	// FindAllStringの第二引数に負の値を指定するとマッチ数の上限が無限になる。
	matches := exp.FindAllString(src, -1)

	result := src
	for _, m := range matches {
		v := NewVariable(m)
		if v == nil {
			continue
		}

		val, err := v.Expand()
		if err != nil {
			return ``, err
		}
		result = strings.Replace(result, m, val, -1)
	}

	return result, nil
}

// 変数名を解析してvariable構造体を生成する。
func NewVariable(key string) *variable {
	if len(key) < minKeyLength {
		return nil
	}

	v := new(variable)
	v.key = key
	v.Place = key[1]
	v.Kind = key[2]

	if v.Place != plcMaster && v.Place != plcServant {
		return nil
	}

	if v.Kind != kndSys && v.Kind != kndEnv && v.Kind != kndJob && v.Kind != kndTime {
		return nil
	}

	fullName := key[3 : len(key)-1]
	nameAndTag := strings.Split(fullName, tagSeparator)
	switch len(nameAndTag) {
	case 1:
		v.Name = nameAndTag[0]
	case 2:
		v.Name = nameAndTag[0]
		v.Tag = nameAndTag[1]
	default:
		return nil
	}

	return v
}

// 変数名の文字列表現を返す。
func (v *variable) String() string {
	return v.key
}

// 変数を値に展開する。
func (v *variable) Expand() (string, error) {
	switch v.Kind {
	case kndSys:
		return v.expandSys()
	case kndEnv:
		return v.expandEnv()
	case kndJob:
		return v.expandJob()
	case kndTime:
		return v.expandTime()
	}

	return ``, fmt.Errorf("Undefined variable[%s].", v)
}

func (v *variable) expandSys() (string, error) {
	fullName := fmt.Sprintf("%s%s%s", v.Name, tagSeparator, v.Tag)
	val, ok := sysValues[fullName]
	if !ok {
		return ``, fmt.Errorf("Undefined variable[%s].", v)
	}

	return val, nil
}

func (v *variable) expandEnv() (string, error) {
	return os.Getenv(v.Name), nil
}

func (v *variable) expandJob() (string, error) {
	if v.Place == plcServant {
		return ``, fmt.Errorf("Cannot use job variable in servant.")
	}
	j, ok := jobValues[v.Name]
	if !ok {
		return ``, fmt.Errorf("Job[%s] is not executed yet.", v.Name)
	}

	switch v.Tag {
	case `ID`:
		return j.ID, nil
	case `RC`:
		return j.RC, nil
	case `SD`:
		t, err := utctime.Parse(utctime.Default, j.SD)
		if err != nil {
			return ``, fmt.Errorf("Cannot parse time string[%s].", j.SD)
		}
		return t.Format("$ST" + utctime.NoDelimiter + "$"), nil
	case `ED`:
		t, err := utctime.Parse(utctime.Default, j.ED)
		if err != nil {
			return ``, fmt.Errorf("Cannot parse time string[%s].", j.ED)
		}
		return t.Format("$ST" + utctime.NoDelimiter + "$"), nil
	case `OUT`:
		return j.OUT, nil
	}

	return ``, fmt.Errorf("Undefined variable[%s].", v)
}

func (v *variable) expandTime() (string, error) {
	if v.Place == plcMaster {
		return ``, fmt.Errorf("Cannot use time variable in master.")
	}

	t, err := utctime.Parse(utctime.NoDelimiter, v.Name)
	if err != nil {
		return ``, fmt.Errorf("Cannot parse time variable[%s]. Reason[%s]", v.Name, err)
	}

	return t.FormatLocaltime(utctime.Default), nil
}
