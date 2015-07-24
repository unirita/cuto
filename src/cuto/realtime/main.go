package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cuto/master/config"
	"cuto/realtime/network"
	"cuto/util"
)

// Runtime arguments
type arguments struct {
	realtimeName string
	jsonMessage  string
}

const usage = `Usage :
    realtime [-n name] json-message

Option :
    -n name : Use realtime network name.

Copyright 2015 unirita Inc.
`

func main() {
	os.Exit(realMain())
}

func realMain() int {
	args := fetchArgs()
	if args == nil {
		showUsage()
		return 1
	}

	configPath := filepath.Join(util.GetRootPath(), "bin", "master.ini")
	if err := config.Load(configPath); err != nil {
		fmt.Println("master.ini not found or cannot read it.")
		return 1
	}
	networkDir := config.Dir.JobnetDir

	if err := network.LoadJobex(args.realtimeName, networkDir); err != nil {
		fmt.Printf("Jobex csv load error: %s\n", err)
	}

	nwk, err := network.Parse(args.jsonMessage)
	if err != nil {
		fmt.Printf("Parse error: %s\n", err)
		return 1
	}

	cmd := network.NewCommand(args.realtimeName)
	nwk.Export(cmd.GetNetworkName(), networkDir)

	id, err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}

	return 0
}

func fetchArgs() *arguments {
	args := new(arguments)
	flag.Usage = showUsage
	flag.StringVar(&args.realtimeName, "n", "", "realtime network name.")
	flag.Parse()
	if flag.NArg() != 1 {
		return nil
	}
	args.jsonMessage = flag.Arg(0)
	return args
}

func showUsage() {
	fmt.Println(usage)
}
