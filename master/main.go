// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/unirita/cuto/console"
	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/db/query"
	"github.com/unirita/cuto/log"
	"github.com/unirita/cuto/master/config"
	"github.com/unirita/cuto/master/jobnet"
	"github.com/unirita/cuto/message"
)

// 実行時引数のオプション
type arguments struct {
	versionFlag   bool   // バージョン情報表示フラグ
	networkName   string // ジョブネットワーク名
	startFlag     bool   // 実行フラグ
	rerunInstance int    // リランを行うインスタンスID
	configPath    string // 設定ファイルのパス
}

// masterの戻り値
const (
	rc_OK    = 0
	rc_ERROR = 1
)

// フラグ系実行時引数のON/OFF
const (
	flag_ON  = true
	flag_OFF = false
)

const defaultConfig = `master.ini`

func main() {
	args := fetchArgs()
	rc := realMain(args)
	os.Exit(rc)
}

func realMain(args *arguments) (rc int) {
	rc = rc_OK

	if args.versionFlag == flag_ON {
		showVersion()
		return
	}

	if args.networkName == "" && args.rerunInstance == 0 {
		showUsage()
		rc = rc_ERROR
		return
	}

	if args.networkName != "" && args.rerunInstance != 0 {
		console.Display("CTM019E", "Cannot use both -n and -r option.")
		rc = rc_ERROR
		return
	}

	if args.configPath == "" {
		args.configPath = defaultConfig
	}

	message.MasterVersion = Version

	if err := config.Load(args.configPath); err != nil {
		console.Display("CTM019E", err)
		console.Display("CTM004E", args.configPath)
		rc = rc_ERROR
		return
	}

	if err := config.DetectError(); err != nil {
		console.Display("CTM005E", err)
		rc = rc_ERROR
		return
	}

	if err := log.Init(config.Dir.LogDir,
		"master",
		"",
		config.Log.OutputLevel,
		config.Log.MaxSizeKB,
		config.Log.MaxGeneration,
		config.Log.TimeoutSec); err != nil {
		console.Display("CTM021E", err)
		rc = rc_ERROR
		return
	}
	defer log.Term()
	console.Display("CTM001I", os.Getpid(), Version)
	// master終了時のコンソール出力
	defer func() {
		console.Display("CTM002I", rc)
	}()

	if args.rerunInstance != 0 {
		nwkResult, err := getNetworkResult(args.rerunInstance)
		if err != nil {
			console.Display("CTM019E", err)
			rc = rc_ERROR
			return
		}

		if nwkResult.Status == db.NORMAL || nwkResult.Status == db.WARN {
			console.Display("CTM029I", args.rerunInstance)
			return
		}

		args.networkName = nwkResult.JobnetWork
		args.startFlag = flag_ON
	}

	nwk := jobnet.LoadNetwork(args.networkName)
	if nwk == nil {
		rc = rc_ERROR
		return
	}
	defer nwk.Terminate()

	if err := nwk.DetectFlowError(); err != nil {
		console.Display("CTM011E", nwk.MasterPath, err)
		rc = rc_ERROR
		return
	}

	if args.startFlag == flag_OFF {
		console.Display("CTM020I", nwk.MasterPath)
		return
	}

	if err := nwk.LoadJobEx(); err != nil {
		console.Display("CTM004E", nwk.JobExPath)
		log.Error(err)
		rc = rc_ERROR
		return
	}

	var err error
	if args.rerunInstance == 0 {
		err = nwk.Run()
	} else {
		nwk.ID = args.rerunInstance
		err = nwk.Rerun()
	}
	if err != nil {
		console.Display("CTM013I", nwk.Name, nwk.ID, "ABNORMAL")
		// ジョブ自体の異常終了では、エラーメッセージが空で返るので、出力しない
		if len(err.Error()) != 0 {
			log.Error(err)
		}
		rc = rc_ERROR
		return
	}
	console.Display("CTM013I", nwk.Name, nwk.ID, "NORMAL")
	return
}

// コマンドライン引数を解析し、arguments構造体を返す。
func fetchArgs() *arguments {
	args := new(arguments)
	flag.Usage = showUsage
	flag.BoolVar(&args.versionFlag, "v", false, "version option")
	flag.StringVar(&args.networkName, "n", "", "network name option")
	flag.BoolVar(&args.startFlag, "s", false, "start option")
	flag.IntVar(&args.rerunInstance, "r", 0, "rerun option")
	flag.StringVar(&args.configPath, "c", "", "config file option")
	flag.Parse()
	return args
}

// バージョンを表示する。
func showVersion() {
	fmt.Printf("%s\n", Version)
}

// オンラインヘルプを表示する。
func showUsage() {
	console.Display("CTM003E")
	fmt.Print(console.USAGE)
}

func getNetworkResult(instanceID int) (*db.JobNetworkResult, error) {
	conn, err := db.Open(config.DB.DBFile)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return query.GetJobnetwork(conn, instanceID)
}
