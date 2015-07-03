package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuto/flowgen/converter"
)

const usage = `Usage :
    flowgen [file_name]

Copyright 2015 unirita Inc.
`

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if len(os.Args) != 2 {
		fmt.Println(usage)
		return 1
	}

	path := os.Args[1]
	elm, err := converter.ParseFile(path)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	definitions := converter.GenerateDefinitions(elm)
	err = converter.ExportFile(convertExtension(path, ".bpmn"), definitions)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
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
