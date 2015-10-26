package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unirita/cuto/flowgen/converter"
)

const usage = `Usage :
    flowgen [file_name]

Copyright 2015 unirita Inc.
`

const (
	rc_OK           = 0
	rc_PARAM_ERROR  = 1
	rc_SYNTAX_ERROR = 2
	rc_OUTPUT_ERROR = 4
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		return rc_PARAM_ERROR
	}

	path := os.Args[1]
	elm, err := converter.ParseFile(path)
	if err != nil {
		fmt.Println(err)
		return rc_SYNTAX_ERROR
	}

	definitions := converter.GenerateDefinitions(elm)
	err = converter.ExportFile(convertExtension(path, ".bpmn"), definitions)
	if err != nil {
		fmt.Println(err)
		return rc_OUTPUT_ERROR
	}

	return rc_OK
}

func convertExtension(path string, afterExt string) string {
	dir, base := filepath.Split(path)
	ext := filepath.Ext(path)

	if ext != "" {
		base = strings.Replace(base, ext, "", -1)
	}
	base = base + afterExt

	return filepath.Join(dir, base)
}
