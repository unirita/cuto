// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package db

import (
	"testing"
)

func TestNewJobResult_初期化できる(t *testing.T) {
	jobNetID := 1

	jobRes := NewJobResult(jobNetID)
	if jobRes.ID != 1 {
		t.Errorf("ジョブネットIDを[%v]で初期化しましたが、[%v]になってしまった。", jobNetID, jobRes.ID)
	}
}
