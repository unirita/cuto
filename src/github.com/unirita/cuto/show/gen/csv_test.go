// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package gen

import (
	"bufio"
	"os"

	"testing"
)

const output_csvfile string = "showtest.csv"

func CreateTestData() *OutputRoot {
	var jobs []*OutputJob

	job := &OutputJob{}
	job.JobId = "job1"
	job.Jobname = "jobName1"
	job.StartDate = "2015-04-27 14:15:24.999"
	job.EndDate = "2015-04-27 14:25:24.999"
	job.Status = 1
	job.Detail = "NORMAL"
	job.Rc = 0
	job.Node = "localhost"
	job.Port = 2015
	job.Variable = "Var1"
	job.CreateDate = "2015-04-27 14:15:24.999"
	job.UpdateDate = "2015-04-27 14:25:24.999"
	jobs = append(jobs, job)

	job = &OutputJob{}
	job.JobId = "job2"
	job.Jobname = "jobName2"
	job.StartDate = "2015-04-27 14:15:24.999"
	job.EndDate = "2015-04-27 14:25:24.999"
	job.Status = 2
	job.Detail = "ABNORMAL"
	job.Rc = 1
	job.Node = "localhost"
	job.Port = 2015
	job.Variable = "Var2"
	job.CreateDate = "2015-04-27 14:15:24.999"
	job.UpdateDate = "2015-04-27 14:25:24.999"
	jobs = append(jobs, job)

	jnet := &OutputJobNet{}
	jnet.Id = 101
	jnet.Jobnetwork = "jn101"
	jnet.StartDate = "2015-04-27 14:15:24.999"
	jnet.EndDate = "2015-04-27 14:25:24.999"
	jnet.Detail = ""
	jnet.Status = 1
	jnet.CreateDate = "2015-04-27 14:15:24.999"
	jnet.UpdateDate = "2015-04-27 14:25:24.999"
	jnet.Jobs = jobs

	root := &OutputRoot{}
	root.Jobnetworks = append(root.Jobnetworks, jnet)
	return root
}

func TestGenerate_CSV形式にジェネレート(t *testing.T) {
	d := CreateTestData()

	var gen CsvGenerator
	msg, err := gen.Generate(d)
	if err != nil {
		t.Fatalf("エラーが返りました。 - %v", err)
	}
	if _, exist := os.Stat(output_csvfile); exist == nil {
		os.Remove(output_csvfile)
	}
	file, err := os.OpenFile(output_csvfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatalf("テスト出力ファイルの生成に失敗しました。 - %v", err)
	}
	_, err = file.WriteString(msg)
	if err != nil {
		t.Fatalf("テスト出力に失敗しました。 - %v", err)
	}
	file.Close()
	file, _ = os.Open(output_csvfile)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	var i int = 0
	for scanner.Scan() {
		line = scanner.Text()
		if i == 0 {
			if line != "#Type,JobNetwork ID,JobNetwork Name,Start Date,End Date,Status,Detail Message,Create Date,Update Date" {
				t.Errorf("不正な行です。[%v]", line)
			}
		} else if i == 1 {
			if line != "JOBNET,101,jn101,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999,1,,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999" {
				t.Errorf("不正な行です。[%v]", line)
			}

		} else if i == 2 {
			if line != "#Type,JobNework ID,Job ID,Job Name,Start Date,End Date,Status,Detail Message,Return Code,Node,Port,Variable,CreateDate,Update Date" {
				t.Errorf("不正な行です。[%v]", line)
			}
		} else if i == 3 {
			if line != "JOB,101,job1,jobName1,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999,1,NORMAL,0,localhost,2015,Var1,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999" {
				t.Errorf("不正な行です。[%v]", line)
			}
		} else if i == 4 {
			if line != "JOB,101,job2,jobName2,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999,2,ABNORMAL,1,localhost,2015,Var2,2015-04-27 14:15:24.999,2015-04-27 14:25:24.999" {
				t.Errorf("不正な行です。[%v]", line)
			}
		}
		i++
	}
}
