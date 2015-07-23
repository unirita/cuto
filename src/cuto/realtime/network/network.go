package network

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	scan "github.com/mattn/go-scan"
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

var jobex = make([][]string, 0)

func init() {
	// Create empty jobex with title record.
	jobex := make([][]string, 1)
	jobex[0] = make([]string, columns)
}

// LoadJobex loads jobex csv which corresponds to name.
// LoadJobex returns empty jobex array if csv is not exists.
func LoadJobex(name string, nwkDir string) error {
	csvPath := searchJobexCsvFile(name, nwkDir)
	if csvPath == "" {
		return nil
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	jobex, err = loadJobexFromReader(file)
	return err
}

func searchJobexCsvFile(name string, nwkDir string) string {
	individualPath := filepath.Join(nwkDir, "realtime", name+".csv")
	defaultPath := filepath.Join(nwkDir, "realtime", "default.csv")

	if _, err := os.Stat(individualPath); !os.IsNotExist(err) {
		return individualPath
	}
	if _, err := os.Stat(defaultPath); !os.IsNotExist(err) {
		return defaultPath
	}

	return ""
}

// loadJobexFromReader reads reader as csv format, and create jobex data array.
func loadJobexFromReader(reader io.Reader) ([][]string, error) {
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

func getJobexRecordByName(jobname string) []string {
	for _, record := range jobex {
		if record[nameIdx] == jobname {
			return record
		}
	}
	return nil
}

type Network struct {
	Flow string `json:"flow"`
	Jobs []Job  `json:"jobs"`
}

// Parse parses str as json format, and create Network object.
func Parse(str string) (*Network, error) {
	decorder := json.NewDecoder(strings.NewReader(str))

	network := new(Network)
	if err := decorder.Decode(network); err != nil {
		return nil, err
	}
	network.complementJobs()

	if err := network.DetectError(); err != nil {
		return nil, err
	}

	return network, nil
}

func (n *Network) complementJobs() {
	for _, record := range jobex {
		isExists := false
		for _, job := range n.Jobs {
			if record[nameIdx] == job.Name {
				isExists = true
				break
			}
		}

		if !isExists {
			newJob := Job{Name: record[nameIdx]}
			newJob.importJobex()
			n.Jobs = append(n.Jobs, newJob)
		}
	}
}

// DetectError detects error in Network object, and return it.
// If there is no error, DetectError returns nil.
func (n *Network) DetectError() error {
	for _, job := range n.Jobs {
		if job.Name == "" {
			return errors.New("Anonymous job detected.")
		}
	}
	return nil
}

type Job struct {
	Name    string
	Node    string
	Port    int
	Path    string
	Param   string
	Env     string
	Work    string
	WRC     int
	WPtn    string
	ERC     int
	EPtn    string
	Timeout int
	SNode   string
	SPort   int
}

// UnmarshalJSON create job object from data(JSON format).
// Use jobex value loaded by LoadJobex function if the parameter is null.
func (j *Job) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if err := scan.ScanTree(i, "/name", &j.Name); err != nil {
		return err
	}
	j.importJobex()

	// scan.ScanTree does not change value of 3rd parameter when error occured.
	scan.ScanTree(i, "/node", &j.Node)
	scan.ScanTree(i, "/port", &j.Port)
	scan.ScanTree(i, "/path", &j.Path)
	scan.ScanTree(i, "/param", &j.Param)
	scan.ScanTree(i, "/env", &j.Env)
	scan.ScanTree(i, "/work", &j.Work)
	scan.ScanTree(i, "/wrc", &j.WRC)
	scan.ScanTree(i, "/wptn", &j.WPtn)
	scan.ScanTree(i, "/erc", &j.ERC)
	scan.ScanTree(i, "/eptn", &j.EPtn)
	scan.ScanTree(i, "/timeout", &j.Timeout)
	scan.ScanTree(i, "/snode", &j.SNode)
	scan.ScanTree(i, "/sport", &j.SPort)

	return nil
}

func (j *Job) importJobex() error {
	for _, record := range jobex {
		if record[nameIdx] == j.Name {
			var err error
			j.Node = record[nodeIdx]
			j.Port, err = strconv.Atoi(record[portIdx])
			if err != nil {
				return err
			}
			j.Path = record[pathIdx]
			j.Param = record[paramIdx]
			j.Env = record[envIdx]
			j.Work = record[workIdx]
			j.WRC, err = strconv.Atoi(record[wrcIdx])
			if err != nil {
				return err
			}
			j.WPtn = record[wptnIdx]
			j.ERC, err = strconv.Atoi(record[ercIdx])
			if err != nil {
				return err
			}
			j.EPtn = record[eptnIdx]
			j.Timeout, err = strconv.Atoi(record[timeoutIdx])
			if err != nil {
				return err
			}
			j.SNode = record[snodeIdx]
			j.SPort, err = strconv.Atoi(record[sportIdx])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
