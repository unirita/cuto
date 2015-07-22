package jobex

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"cuto/realtime/network"
)

// Number of columns
const columns = 14

// Indexes of columns
const (
	nameIdx = iota
	nodeIdx
	portIdx
	pathIdx
	paramIdx
	envIdx
	workIdx
	wrcIdx
	wptnIdx
	ercIdx
	eptnIdx
	timeoutIdx
	snodeIdx
	sportIdx
)

// Load reads reader as csv format, and create jobex data array.
func Load(reader io.Reader) ([][]string, error) {
	r := csv.NewReader(reader)
	jobex, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(jobex) > 0 && len(jobex[0]) != columns {
		return nil, fmt.Errorf("Number of jobex csv columns[%d] must be %d.", len(jobex[0]), columns)
	}

	return jobex, nil
}

// MergeRealtimeParam merge jobex data from jobs to base.
func MergeRealtimeParam(base [][]string, jobs []network.Job) [][]string {
	result := make([][]string, len(base))
	copy(result, base)

	for _, job := range jobs {
		isExistsInBase := false
		for idx, baseRecord := range base {
			if baseRecord[nameIdx] == *job.Name {
				isExistsInBase = true
				result[idx] = mergeRecord(baseRecord, job)
				break
			}
		}

		if !isExistsInBase {
			newRecord := make([]string, columns)
			newRecord[nameIdx] = *job.Name
			result = append(result, mergeRecord(newRecord, job))
		}
	}
	return result
}

func mergeRecord(record []string, job network.Job) []string {
	result := make([]string, columns)
	for idx, col := range record {
		switch idx {
		case nameIdx:
			result[nameIdx] = col
		case nodeIdx:
			if job.Node != nil {
				result[nodeIdx] = *job.Node
			} else {
				result[nodeIdx] = col
			}
		case portIdx:
			if job.Port != nil {
				result[portIdx] = strconv.Itoa(*job.Port)
			} else {
				result[portIdx] = col
			}
		case pathIdx:
			if job.Path != nil {
				result[pathIdx] = *job.Path
			} else {
				result[pathIdx] = col
			}
		case paramIdx:
			if job.Param != nil {
				result[paramIdx] = *job.Param
			} else {
				result[paramIdx] = col
			}
		case envIdx:
			if job.Env != nil {
				result[envIdx] = *job.Env
			} else {
				result[envIdx] = col
			}
		case workIdx:
			if job.Work != nil {
				result[workIdx] = *job.Work
			} else {
				result[workIdx] = col
			}
		case wrcIdx:
			if job.WRC != nil {
				result[wrcIdx] = strconv.Itoa(*job.WRC)
			} else {
				result[wrcIdx] = col
			}
		case wptnIdx:
			if job.WPtn != nil {
				result[wptnIdx] = *job.WPtn
			} else {
				result[wptnIdx] = col
			}
		case ercIdx:
			if job.ERC != nil {
				result[ercIdx] = strconv.Itoa(*job.ERC)
			} else {
				result[ercIdx] = col
			}
		case eptnIdx:
			if job.EPtn != nil {
				result[eptnIdx] = *job.EPtn
			} else {
				result[eptnIdx] = col
			}
		case timeoutIdx:
			if job.Timeout != nil {
				result[timeoutIdx] = strconv.Itoa(*job.Timeout)
			} else {
				result[timeoutIdx] = col
			}
		case snodeIdx:
			if job.SNode != nil {
				result[snodeIdx] = *job.SNode
			} else {
				result[snodeIdx] = col
			}
		case sportIdx:
			if job.SPort != nil {
				result[sportIdx] = strconv.Itoa(*job.SPort)
			} else {
				result[sportIdx] = col
			}
		}
	}

	return result
}
