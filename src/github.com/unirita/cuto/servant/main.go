// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"cuto/console"
	"cuto/log"
	"cuto/message"
	"cuto/util"

	"cuto/servant/config"
)

type arguments struct {
	v          bool
	configPath string // 設定ファイルのパス
}

const (
	rc_OK    = 0
	rc_error = 1
)

// エントリポイント
func main() {
	args := fetchArgs()
	// 作業フォルダは実行モジュールと同じ場所にする
	os.Chdir(util.GetCurrentPath())

	os.Exit(realMain(args))
}

func realMain(args *arguments) int {
	if args.v {
		showVersion()
		return rc_OK
	}
	message.ServantVersion = Version

	// システム変数のセット
	message.AddSysValue("ROOT", "", util.GetRootPath())

	config.ReadConfig(args.configPath)
	if err := config.Servant.DetectError(); err != nil {
		console.Display("CTS005E", err)
		return rc_error
	}

	// ログ出力開始
	if err := log.Init(config.Servant.Dir.LogDir,
		"servant",
		strconv.Itoa(config.Servant.Sys.BindPort),
		config.Servant.Log.OutputLevel,
		config.Servant.Log.MaxSizeKB,
		config.Servant.Log.MaxGeneration,
		config.Servant.Log.TimeoutSec); err != nil {
		console.Display("CTS023E", err)
		return rc_error
	}
	defer log.Term()
	console.Display("CTS001I", os.Getpid(), Version)

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
	flag.StringVar(&args.configPath, "c", "", "config file option")
	flag.Parse()
	return args
}

func showVersion() {
	fmt.Printf("cuto servant version %s\n", Version)
}
