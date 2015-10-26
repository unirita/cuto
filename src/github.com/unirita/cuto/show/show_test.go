// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package main

import (
	"testing"

	"cuto/db"
	"cuto/show/gen"
)

func TestNewShowParam_ShowParamの初期化(t *testing.T) {
	var g gen.Generator
	var c gen.CsvGenerator
	g = c
	s := NewShowParam(123, "jobnet1", "2015-04-27 14:15:24.999", "2015-04-27 15:15:24.999", 9, g)

	if s.nid != 123 {
		t.Errorf("不正なインスタンスID[%v]が返りました。", s.nid)
	}
	if s.jobnetName != "jobnet1" {
		t.Errorf("不正なジョブネット名[%v]が返りました。", s.jobnetName)
	}
	if s.from != "2015-04-27 14:15:24.999" {
		t.Errorf("不正なfrom[%v]が返りました。", s.from)
	}
	if s.to != "2015-04-27 15:15:24.999" {
		t.Errorf("to[%v]が返りました。", s.to)
	}
	if s.status != 9 {
		t.Errorf("不正なstatus[%v]が返りました。", s.status)
	}
	switch s.gen.(type) {
	case gen.CsvGenerator:
	default:
		t.Error("不正なジェネレーターです。")
	}
}

func TestSetOutputStructure_出力情報のセッティング(t *testing.T) {
	// *** テストデータ作成 ***
	jobnet := &oneJobnetwork{}
	jnetRes := &db.JobNetworkResult{
		ID:         123,
		JobnetWork: "Jobnet123",
		StartDate:  "2015-04-27 14:15:24.999",
		EndDate:    "2015-04-27 15:15:24.999",
		Status:     1,
		Detail:     "",
		CreateDate: "2015-04-27 14:15:24.999",
		UpdateDate: "2015-04-27 15:15:24.999",
	}
	jobnet.jobnet = jnetRes
	jRes1 := db.JobResult{
		ID:         123,
		JobId:      "Job1",
		JobName:    "Jobname1",
		StartDate:  "2015-04-27 14:15:24.999",
		EndDate:    "2015-04-27 14:25:24.999",
		Status:     1,
		Detail:     "d1",
		Rc:         0,
		Node:       "localhost",
		Port:       2015,
		Variable:   "Var123",
		CreateDate: "2015-04-27 14:15:24.999",
		UpdateDate: "2015-04-27 14:25:24.999",
	}
	jobnet.jobs = append(jobnet.jobs, &jRes1)

	jRes2 := db.JobResult{
		ID:         124,
		JobId:      "Job2",
		JobName:    "Jobname2",
		StartDate:  "2015-04-27 14:25:24.999",
		EndDate:    "2015-04-27 14:35:24.999",
		Status:     1,
		Detail:     "d2",
		Rc:         0,
		Node:       "localhost",
		Port:       2015,
		Variable:   "Var456",
		CreateDate: "2015-04-27 14:25:24.999",
		UpdateDate: "2015-04-27 14:35:24.999",
	}
	jobnet.jobs = append(jobnet.jobs, &jRes2)
	// *** ここまで ***
	output := jobnet.setOutputStructure(true)
	if output.Id != jnetRes.ID {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.ID, output.Id)
	}
	if output.Jobnetwork != jnetRes.JobnetWork {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.JobnetWork, output.Jobnetwork)
	}
	if output.StartDate != jnetRes.StartDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.StartDate, output.StartDate)
	}
	if output.EndDate != jnetRes.EndDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.EndDate, output.EndDate)
	}
	if output.Status != jnetRes.Status {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.Status, output.Status)
	}
	if len(output.Detail) != 0 {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.Detail, output.Detail)
	}
	if output.CreateDate != jnetRes.CreateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.CreateDate, output.CreateDate)
	}
	if output.UpdateDate != jnetRes.UpdateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jnetRes.UpdateDate, output.UpdateDate)
	}
	if output.Jobs[0].JobId != jRes1.JobId {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.JobId, output.Jobs[0].JobId)
	}
	if output.Jobs[0].Jobname != jRes1.JobName {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.JobName, output.Jobs[0].Jobname)
	}
	if output.Jobs[0].StartDate != jRes1.StartDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.StartDate, output.Jobs[0].StartDate)
	}
	if output.Jobs[0].EndDate != jRes1.EndDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.EndDate, output.Jobs[0].EndDate)
	}
	if output.Jobs[0].Detail != jRes1.Detail {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Detail, output.Jobs[0].Detail)
	}
	if output.Jobs[0].Rc != jRes1.Rc {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Rc, output.Jobs[0].Rc)
	}
	if output.Jobs[0].Status != jRes1.Status {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Status, output.Jobs[0].Status)
	}
	if output.Jobs[0].Node != jRes1.Node {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Node, output.Jobs[0].Node)
	}
	if output.Jobs[0].Port != jRes1.Port {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Port, output.Jobs[0].Port)
	}
	if output.Jobs[0].Variable != jRes1.Variable {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.Variable, output.Jobs[0].Variable)
	}
	if output.Jobs[0].CreateDate != jRes1.CreateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.CreateDate, output.Jobs[0].CreateDate)
	}
	if output.Jobs[0].UpdateDate != jRes1.UpdateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes1.UpdateDate, output.Jobs[0].UpdateDate)
	}

	if output.Jobs[1].JobId != jRes2.JobId {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.JobId, output.Jobs[1].JobId)
	}
	if output.Jobs[1].Jobname != jRes2.JobName {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.JobName, output.Jobs[1].Jobname)
	}
	if output.Jobs[1].StartDate != jRes2.StartDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.StartDate, output.Jobs[1].StartDate)
	}
	if output.Jobs[1].EndDate != jRes2.EndDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.EndDate, output.Jobs[1].EndDate)
	}
	if output.Jobs[1].Detail != jRes2.Detail {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Detail, output.Jobs[1].Detail)
	}
	if output.Jobs[1].Rc != jRes2.Rc {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Rc, output.Jobs[1].Rc)
	}
	if output.Jobs[1].Status != jRes2.Status {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Status, output.Jobs[1].Status)
	}
	if output.Jobs[1].Node != jRes2.Node {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Node, output.Jobs[1].Node)
	}
	if output.Jobs[1].Port != jRes2.Port {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Port, output.Jobs[1].Port)
	}
	if output.Jobs[1].Variable != jRes2.Variable {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.Variable, output.Jobs[1].Variable)
	}
	if output.Jobs[1].CreateDate != jRes2.CreateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.CreateDate, output.Jobs[1].CreateDate)
	}
	if output.Jobs[1].UpdateDate != jRes2.UpdateDate {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", jRes2.UpdateDate, output.Jobs[1].UpdateDate)
	}

}

func TestSetOutputStructure_ローカルタイムゾーンでの出力(t *testing.T) {
	// *** テストデータ作成 ***
	jobnet := &oneJobnetwork{}
	jnetRes := &db.JobNetworkResult{
		ID:         123,
		JobnetWork: "Jobnet123",
		StartDate:  "2015-04-27 14:15:24.999",
		EndDate:    "2015-04-27 15:15:24.999",
		Status:     1,
		Detail:     "",
		CreateDate: "2015-04-27 14:15:24.999",
		UpdateDate: "2015-04-27 15:15:24.999",
	}
	jobnet.jobnet = jnetRes
	jRes := db.JobResult{
		ID:         123,
		JobId:      "Job1",
		JobName:    "Jobname1",
		StartDate:  "2015-04-27 14:15:24.999",
		EndDate:    "2015-04-27 14:25:24.999",
		Status:     1,
		Detail:     "d1",
		Rc:         0,
		Node:       "localhost",
		Port:       2015,
		Variable:   "Var123",
		CreateDate: "2015-04-27 14:15:24.999",
		UpdateDate: "2015-04-27 14:25:24.999",
	}
	jobnet.jobs = append(jobnet.jobs, &jRes)
	// *** ここまで ***
	output := jobnet.setOutputStructure(false)
	if output.StartDate != "2015-04-27 23:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:15:24.999", output.StartDate)
	}
	if output.EndDate != "2015-04-28 00:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-28 00:15:24.999", output.EndDate)
	}
	if output.CreateDate != "2015-04-27 23:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:15:24.999", output.CreateDate)
	}
	if output.UpdateDate != "2015-04-28 00:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-28 00:15:24.999", output.UpdateDate)
	}
	if output.Jobs[0].StartDate != "2015-04-27 23:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:15:24.999", output.Jobs[0].StartDate)
	}
	if output.Jobs[0].EndDate != "2015-04-27 23:25:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:25:24.999", output.Jobs[0].EndDate)
	}
	if output.Jobs[0].CreateDate != "2015-04-27 23:15:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:15:24.999", output.Jobs[0].CreateDate)
	}
	if output.Jobs[0].UpdateDate != "2015-04-27 23:25:24.999" {
		t.Errorf("[%v]であるべきデータが、[%v]と返りました。", "2015-04-27 23:25:24.999", output.Jobs[0].UpdateDate)
	}

}
