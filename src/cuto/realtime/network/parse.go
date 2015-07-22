package network

import (
	"encoding/json"
	"errors"
	"io"
)

type Network struct {
	Flow string `json:"flow"`
	Jobs []Job  `json:"jobs"`
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

// Parse reads json string from reader, and unmarshal it.
func Parse(reader io.Reader) (*Network, error) {
	decorder := json.NewDecoder(reader)

	network := new(Network)
	if err := decorder.Decode(network); err != nil {
		return nil, err
	}
	if err := network.DetectError(); err != nil {
		return nil, err
	}

	return network, nil
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
