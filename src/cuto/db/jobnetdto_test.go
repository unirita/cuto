// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

import (
	"testing"

	"cuto/utctime"
)

func TestNewJobNetworkResult_初期化できる(t *testing.T) {
	name := "ABC"
	start := utctime.Now().String()
	status := NORMAL

	net := NewJobNetworkResult(name, start, status)
	if net.JobnetWork != name {
		t.Errorf("ジョブネット名を[%v]で初期化しましたが、[%v]が返りました。", name, net.JobnetWork)
	}
	if net.StartDate != start {
		t.Errorf("開始日時を[%v]で初期化しましたが、[%v]が返りました。", start, net.StartDate)
	}
	if net.Status != status {
		t.Errorf("ジョブネット名を[%v]で初期化しましたが、[%v]が返りました。", status, net.Status)
	}
}
