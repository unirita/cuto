// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package gen

import (
	"testing"
)

func TestGenerate_JSON形式にジェネレート(t *testing.T) {
	d := CreateTestData()

	var gen JsonGenerator
	msg, err := gen.Generate(d)
	if err != nil {
		t.Fatalf("エラーが返りました。 - %v", err)
	}
	if msg != "{\"jobnetworks\":[{\"id\":101,\"jobnetwork\":\"jn101\",\"startdate\":\"2015-04-27 14:15:24.999\",\"enddate\":\"2015-04-27 14:25:24.999\",\"status\":1,\"detail\":\"\",\"createdate\":\"2015-04-27 14:15:24.999\",\"updatedate\":\"2015-04-27 14:25:24.999\",\"jobs\":[{\"jobid\":\"job1\",\"jobname\":\"jobName1\",\"startdate\":\"2015-04-27 14:15:24.999\",\"enddate\":\"2015-04-27 14:25:24.999\",\"status\":1,\"detail\":\"NORMAL\",\"rc\":0,\"Node\":\"localhost\",\"port\":2015,\"variable\":\"Var1\",\"createdate\":\"2015-04-27 14:15:24.999\",\"updatedate\":\"2015-04-27 14:25:24.999\"},{\"jobid\":\"job2\",\"jobname\":\"jobName2\",\"startdate\":\"2015-04-27 14:15:24.999\",\"enddate\":\"2015-04-27 14:25:24.999\",\"status\":2,\"detail\":\"ABNORMAL\",\"rc\":1,\"Node\":\"localhost\",\"port\":2015,\"variable\":\"Var2\",\"createdate\":\"2015-04-27 14:15:24.999\",\"updatedate\":\"2015-04-27 14:25:24.999\"}]}]}" {
		t.Errorf("不正なデータが返りました。 - %v", msg)
	}
}
