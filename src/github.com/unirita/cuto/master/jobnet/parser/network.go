// ジョブネットワーク定義ファイル全体の定義。
// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// ネットワーク定義BPMNのdefinitions要素。
type Definitions struct {
	Process []Process `xml:"process"`
}

// ネットワーク定義BPMNのprocess要素。
type Process struct {
	Start   []StartEvent      `xml:"startEvent"`
	End     []EndEvent        `xml:"endEvent"`
	Task    []ServiceTask     `xml:"serviceTask"`
	Gateway []ParallelGateway `xml:"parallelGateway"`
	Flow    []SequenceFlow    `xml:"sequenceFlow"`
}

// ネットワーク定義BPMNのstartEvent要素。
type StartEvent struct {
	ID string `xml:"id,attr"`
}

// ネットワーク定義BPMNのendEvent要素。
type EndEvent struct {
	ID string `xml:"id,attr"`
}

// ネットワーク定義BPMNのserviceTask要素。
type ServiceTask struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// ネットワーク定義BPMNのparallelGateway要素。
type ParallelGateway struct {
	ID string `xml:"id,attr"`
}

// ネットワーク定義BPMNのsequenceFlow要素。
type SequenceFlow struct {
	From string `xml:"sourceRef,attr"`
	To   string `xml:"targetRef,attr"`
}

// ネットワーク定義をファイルから読み込み、パース結果を返す
func ParseNetworkFile(fileName string) (*Process, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseNetwork(file)
}

// ネットワーク定義をio.Readerから読み込み、パース結果を返す
func ParseNetwork(reader io.Reader) (*Process, error) {
	dec := xml.NewDecoder(reader)

	defs := new(Definitions)
	err := dec.Decode(defs)
	if err != nil {
		return nil, err
	}

	if len(defs.Process) != 1 {
		return nil, fmt.Errorf("Process element is required, and must be unique.")
	}
	proc := defs.Process[0]

	if len(proc.Start) != 1 {
		return nil, fmt.Errorf("StartEvent element is required, and must be unique.")
	}

	if len(proc.End) != 1 {
		return nil, fmt.Errorf("EndEvent element is required, and must be unique.")
	}

	return &proc, nil
}
