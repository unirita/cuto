package converter

import (
	"errors"
	"regexp"
	"strings"
)

const (
	arrow  = "->"
	gHead  = "["
	gTerm  = "]"
	gDelim = ","
	tmpGW  = ":gw"
)

// ParseString parses flow description from string.
// If it successed, returns head element of flow.
func ParseString(str string) (Element, error) {
	str = ignoreSpaceChars(str)
	extracted, gws, err := extractGateways(str)

	var head Element = nil
	elmStrs := strings.Split(extracted, arrow)
	isFirst := true
	var current Element
	var pre Element
	var tmpGWIndex = 0
	for _, elmStr := range elmStrs {
		if elmStr == tmpGW {
			if tmpGWIndex > len(gws) {
				return nil, errors.New("Gateway parse error.")
			}
			current = gws[tmpGWIndex]
			tmpGWIndex++
		} else {
			current, err = parseJob(elmStr)
			if err != nil {
				return nil, err
			}
		}

		if isFirst {
			head = current
			isFirst = false
		} else {
			pre.SetNext(current)
		}
		pre = current
	}
	return head, nil
}

func ignoreSpaceChars(str string) string {
	ptn := regexp.MustCompile(`\s`)
	return ptn.ReplaceAllLiteralString(str, "")
}

func extractGateways(str string) (string, []*Gateway, error) {
	ptn := regexp.MustCompile(`\[.+?\]`)
	gws := make([]*Gateway, 0)

	gwStrs := ptn.FindAllString(str, -1)
	for _, gwStr := range gwStrs {
		gw, err := parseGateway(gwStr)
		if err != nil {
			return str, nil, err
		}
		gws = append(gws, gw)
	}

	str = ptn.ReplaceAllString(str, tmpGW)
	return str, gws, nil
}

func parseJob(str string) (*Job, error) {
	ptn := regexp.MustCompile(`[\[\]\\/:*?<>$&,\-]`)
	if ptn.MatchString(str) {
		return nil, errors.New("Irregal character found in job name")
	}
	return NewJob(str), nil
}

func parseGateway(str string) (*Gateway, error) {
	str = strings.TrimLeft(str, gHead)
	str = strings.TrimRight(str, gTerm)

	p := NewGateway()

	pathStrs := strings.Split(str, gDelim)
	for _, pathStr := range pathStrs {
		jobNames := strings.Split(pathStr, arrow)
		isFirst := true
		var preJob *Job
		for _, jobName := range jobNames {
			j, err := parseJob(jobName)
			if err != nil {
				return nil, err
			}

			if isFirst {
				p.AddPathHead(j)
				isFirst = false
			} else {
				preJob.SetNext(j)
			}
			preJob = j
		}
	}

	return p, nil
}
