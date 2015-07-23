package network

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

type Network struct {
	Flow string `json:"flow"`
	Jobs []Job  `json:"jobs"`
}

// DetectError detects error in Network object, and return it.
// If there is no error, DetectError returns nil.
func (n *Network) DetectError() error {
	for _, job := range n.Jobs {
		if job.Name == nil || *job.Name == "" {
			return errors.New("Anonymous job detected.")
		}
	}
	return nil
}

type Job struct {
	Name    *string `json:"name"`
	Node    *string `json:"node"`
	Port    *int    `json:"port"`
	Path    *string `json:"path"`
	Param   *string `json:"param"`
	Env     *string `json:"env"`
	Work    *string `json:"work"`
	WRC     *int    `json:"wrc"`
	WPtn    *string `json:"wptn"`
	ERC     *int    `json:"erc"`
	EPtn    *string `json:"eptn"`
	Timeout *int    `json:"timeout"`
	SNode   *string `json:"snode"`
	SPort   *int    `json:"sport"`
}

// Parse parses str as json format, and create Network object.
func Parse(str string) (*Network, error) {
	decorder := json.NewDecoder(strings.NewReader(str))

	network := new(Network)
	if err := decorder.Decode(network); err != nil {
		return nil, err
	}
	if err := network.DetectError(); err != nil {
		return nil, err
	}

	return network, nil
}
