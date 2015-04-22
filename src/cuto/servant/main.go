// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package main

import (
	"flag"
	"fmt"
	"os"

	"cuto/console"
	"cuto/log"
	"cuto/message"
	"cuto/util"

	"cuto/servant/config"
)

type arguments struct {
	v bool
}

const (
	rc_OK    = 0
	rc_error = 1
)

// エントリポイント
func main() {
	os.Exit(realMain())
}

func realMain() int {
	args := fetchArgs()
	if args.v {
		showVersion()
		return rc_OK
	}

	// システム変数のセット
	message.AddSysValue("ROOT", "", util.GetRootPath())

	console.Display("CTS001I", Version)

	config.ReadConfig()
	if err := config.Servant.DetectError(); err != nil {
		console.Display("CTS005E", err)
		return rc_error
	}

	// ログ出力開始
	if err := log.Init(config.Servant.Dir.LogDir,
		"servant",
		config.Servant.Log.OutputLevel,
		config.Servant.Log.MaxSizeKB,
		config.Servant.Log.MaxGeneration); err != nil {
		console.Display("CTS023E", err)
		return rc_error
	}
	defer log.Term()

	// メイン処理開始
	exitCode, err := Run()

	if err != nil {
		log.Error(err)
		exitCode = rc_error
	}
	console.Display("CTS002I", exitCode)
	return exitCode
}

func fetchArgs() *arguments {
	args := new(arguments)
	flag.BoolVar(&args.v, "v", false, "version option")
	flag.Parse()
	return args
}

func showVersion() {
	fmt.Printf("cuto servant version %s\n", Version)
}
