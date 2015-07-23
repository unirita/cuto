package network

import (
	"strconv"
)

// MergeParamIntoJobex merge jobex data from jobs to base.
func MergeParamIntoJobex(base [][]string, jobs []Job) [][]string {
	result := make([][]string, len(base))
	copy(result, base)

	for _, job := range jobs {
		isExistsInBase := false
		for idx, baseRecord := range base {
			if baseRecord[nameIdx] == *job.Name {
				isExistsInBase = true
				result[idx] = mergeJobexRecord(baseRecord, job)
				break
			}
		}

		if !isExistsInBase {
			newRecord := make([]string, columns)
			newRecord[nameIdx] = *job.Name
			result = append(result, mergeJobexRecord(newRecord, job))
		}
	}
	return result
}

func mergeJobexRecord(record []string, job Job) []string {
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
