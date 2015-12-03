// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package gen

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// JSON形式のジェネレーター
type CsvGenerator struct {
}

func (s CsvGenerator) Generate(out *OutputRoot) (string, error) {
	var jnBuf, jobBuf bytes.Buffer

	jnWriter := csv.NewWriter(&jnBuf)
	jnWriter.Write([]string{"#Type", "JobNetwork ID", "JobNetwork Name", "Start Date", "End Date",
		"Status", "Detail Message", "Create Date", "Update Date"})
	jobWriter := csv.NewWriter(&jobBuf)
	jobWriter.Write([]string{"#Type", "JobNework ID", "Job ID", "Job Name", "Start Date", "End Date",
		"Status", "Detail Message", "Return Code", "Node", "Port", "Variable", "CreateDate", "Update Date"})
	for _, jn := range out.Jobnetworks {
		if err := jnWriter.Write([]string{"JOBNET", fmt.Sprintf("%d", jn.Id),
			jn.Jobnetwork, jn.StartDate, jn.EndDate, fmt.Sprintf("%d", jn.Status),
			jn.Detail, jn.CreateDate, jn.UpdateDate}); err != nil {

			panic(err)
		}
		for _, job := range jn.Jobs {
			if err := jobWriter.Write([]string{"JOB", fmt.Sprintf("%d", jn.Id), job.JobId, job.Jobname,
				job.StartDate, job.EndDate, fmt.Sprintf("%d", job.Status), job.Detail,
				fmt.Sprintf("%d", job.Rc), job.Node,
				fmt.Sprintf("%d", job.Port), job.Variable, job.CreateDate, job.UpdateDate}); err != nil {

				panic(err)
			}
		}
	}
	jnWriter.Flush()
	jobWriter.Flush()
	return jnBuf.String() + jobBuf.String(), nil
}
