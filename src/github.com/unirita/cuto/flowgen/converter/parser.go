package converter

import (
	"errors"
	"io/ioutil"
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

// Parse parses flow description from file.
func ParseFile(filepath string) (Element, error) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return ParseString(string(buf))
}

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
	if err := validateJobName(str); err != nil {
		return nil, err
	}
	return NewJob(str), nil
}

func validateJobName(str string) error {
	if len(str) == 0 {
		return errors.New("Empty job name found.")
	}

	ptn := regexp.MustCompile(`[\[\]\\/:*?<>$&,\-]`)
	if ptn.MatchString(str) {
		return errors.New("Irregal character found in job name.")
	}

	return nil
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
