// Copyright 2015 unirita Inc.
// Created 2015/04/14 shanxia

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"path/filepath"

	"cuto/console"
	"cuto/db"
	"cuto/master/config"
	"cuto/show/gen"
	"cuto/utctime"
)

// 実行時引数のオプション
type arguments struct {
	help   bool   // Usageを表示
	v      bool   // バージョン情報表示
	nid    int    // ジョブネットワークのインスタンスID
	jobnet string // ジョブネットワーク名
	from   string // From日付
	to     string // To日付
	status string // ジョブネットワークのステータス
	format string // 表示フォーマット
	config string // 設定ファイルのパス
	isUTC  bool   // 時刻を標準時として扱うかどうか
}

// 戻り値
const (
	rc_OK      = 0  // 正常終了
	rc_NOTHING = 4  // 出力件数が0件
	rc_PARMERR = 8  // パラメータエラー
	rc_ERROR   = 12 // 実行時エラー
)

// デフォルトの設定ファイル名
var defaultConfig string = getDefaultConfig()

func getDefaultConfig() string {
	if runtime.GOOS == "windows" {
		return "master.ini"
	}
	return filepath.Join(os.Getenv("CUTOROOT"), "bin", "master.ini")
}

// ステータス指定
const status_normal = "normal"
const status_abnormal = "abnormal"
const status_running = "running"
const status_warn = "warn"

func main() {
	console.DisplayError("CTU001I", Version)

	rc := realMain(fetchArgs())

	console.DisplayError("CTU002I", rc)
	os.Exit(rc)
}

// 処理のメインルーチン
func realMain(args *arguments) int {
	if args.v { // バージョン情報表示
		showVersion()
		return rc_OK
	}
	if args.help { // Usage表示
		showUsage()
		return rc_OK
	}
	// 設定ファイル名の取得
	if len(args.config) == 0 {
		args.config = defaultConfig
	}
	if err := config.Load(args.config); err != nil { // 設定ファイル読み込み。
		console.DisplayError("CTU006E", args.config)
		return rc_PARMERR
	}
	if len(args.from) == 0 && len(args.to) == 0 { // From-to指定無しの場合は、現在のCPU日付のみを対象とする
		now := utctime.Now()
		if args.isUTC {
			args.from = now.Format(utctime.Date8Num)
			args.to = args.from
		} else {
			args.from = now.FormatLocaltime(utctime.Date8Num)
			args.to = args.from
		}
	}
	from, to, err := parseFromTo(args.from, args.to, args.isUTC)
	if err != nil {
		console.DisplayError("CTU003E", err)
		showUsage()
		return rc_PARMERR
	}
	status, err := getStatusType(args.status) // status取得
	if err != nil {
		console.DisplayError("CTU003E", fmt.Sprintf("Invalid status option. [%v]", args.status))
		showUsage()
		return rc_PARMERR
	}
	gen := getSeparatorType(args.format) // 出力形態
	if gen == nil {
		console.DisplayError("CTU003E", fmt.Sprintf("Invalid [format] format.[%v]", args.format))
		showUsage()
		return rc_PARMERR
	}
	param := NewShowParam(args.nid, args.jobnet, from, to, status, &gen)
	rc, err := param.Run(config.DB.DBFile)
	if err != nil {
		console.DisplayError("CTU004E", err)
		return rc_ERROR
	} else if rc == 0 { // 出力件数が0
		return rc_NOTHING
	}
	return rc_OK
}

// 引数情報の取得
func fetchArgs() *arguments {
	args := new(arguments)
	flag.Usage = showUsage
	flag.BoolVar(&args.help, "help", false, "usage option.")
	flag.BoolVar(&args.v, "v", false, "version option.")
	flag.IntVar(&args.nid, "nid", 0, "Jobnetwork instance ID.")
	flag.StringVar(&args.jobnet, "jobnet", "", "jobnet name option.")
	flag.StringVar(&args.from, "from", "", "From date.")
	flag.StringVar(&args.to, "to", "", "To date.")
	flag.StringVar(&args.status, "status", "", "Jobnetwork status.")
	flag.StringVar(&args.format, "format", "", "Output format.")
	flag.StringVar(&args.config, "c", "", "Input config-file.")
	flag.BoolVar(&args.isUTC, "utc", false, "UTC option.")
	flag.Parse()
	return args
}

// バージョン情報の表示
func showVersion() {
	fmt.Fprintf(os.Stderr, "%v show-utility version.\n", Version)
}

// usage情報の表示
func showUsage() {
	fmt.Fprintf(os.Stderr, console.USAGE_SHOW)
}

// 出力形態の取得
func getSeparatorType(value string) gen.Generator {
	if len(value) == 0 || value == "json" {
		return *new(gen.JsonGenerator)
	} else if value == "csv" {
		return *new(gen.CsvGenerator)
	}
	return nil
}

// 取得するジョブネットステータスの取得
func getStatusType(status string) (int, error) {
	if len(status) == 0 { // ステータス指定無し
		return -1, nil
	} else if status == status_normal {
		return db.NORMAL, nil
	} else if status == status_abnormal {
		return db.ABNORMAL, nil
	} else if status == status_running {
		return db.RUNNING, nil
	} else if status == status_warn {
		return db.WARN, nil
	}
	return 0, fmt.Errorf("Unknown status type [%v].", status)
}

func parseFromTo(fromArg, toArg string, isUTC bool) (string, string, error) {
	parseMethod := utctime.Parse
	if !isUTC {
		parseMethod = utctime.ParseLocaltime
	}
	var from, to string
	if len(fromArg) != 0 {
		f, err := parseMethod(utctime.Date8Num, fromArg)
		if err != nil {
			return "", "", fmt.Errorf("Invalid [from] format. [%s]", fromArg)
		}
		from = f.String()
	}
	if len(toArg) != 0 {
		t, err := parseMethod(utctime.NoDelimiter, toArg+"235959.999")
		if err != nil {
			return "", "", fmt.Errorf("Invalid [to] format. [%s]", toArg)
		}

		to = t.String()
	}

	return from, to, nil
}
