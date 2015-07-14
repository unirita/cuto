// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"cuto/util"
)

type config struct {
	Job jobSection
	Dir dirSection
	DB  dbSection
	Log logSection
}

// 設定ファイルのjobセクション
type jobSection struct {
	DefaultNode          string `toml:"default_node"`
	DefaultPort          int    `toml:"default_port"`
	DefaultTimeoutMin    int    `toml:"default_timeout_min"`
	ConnectionTimeoutSec int    `toml:"connection_timeout_sec"`
	TimeTrackingSpanMin  int    `toml:"time_tracking_span_min"`
	AttemptLimit         int    `toml:"attempt_limit"`
}

// 設定ファイルのdirセクション
type dirSection struct {
	JobnetDir string `toml:"jobnet_dir"`
	LogDir    string `toml:"log_dir"`
}

// 設定ファイルのdbセクション
type dbSection struct {
	DBFile string `toml:"db_file"`
}

// 設定ファイルのlogセクション
type logSection struct {
	OutputLevel   string `toml:"output_level"`
	MaxSizeKB     int    `toml:"max_size_kb"`
	MaxGeneration int    `toml:"max_generation"`
	TimeoutSec    int    `toml:"timeout_sec"`
}

const tag_CUTOROOT = "<CUTOROOT>"

var Dir = new(dirSection)
var Job = new(jobSection)
var DB = new(dbSection)
var Log = new(logSection)

// 設定ファイルをロードする。
//
// 引数: filePath ロードする設定ファイルのパス
//
// 戻り値： エラー情報
func Load(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	return loadReader(f)
}

func loadReader(reader io.Reader) error {
	c := new(config)
	c.Job.AttemptLimit = 1
	if _, err := toml.DecodeReader(reader, c); err != nil {
		return err
	}

	replaceCutoroot(c)

	Dir = &c.Dir
	Job = &c.Job
	DB = &c.DB
	Log = &c.Log
	return nil
}

func replaceCutoroot(c *config) {
	c.Dir.JobnetDir = strings.Replace(c.Dir.JobnetDir, tag_CUTOROOT, util.GetRootPath(), -1)
	c.Dir.LogDir = strings.Replace(c.Dir.LogDir, tag_CUTOROOT, util.GetRootPath(), -1)
}

// 設定値のエラー検出を行う。
//
// return : エラー情報
func DetectError() error {
	if Job.DefaultPort < 0 || 65535 < Job.DefaultPort {
		return fmt.Errorf("job.default_port(%d) must be within the range 0 and 65535.", Job.DefaultPort)
	}
	if Job.DefaultTimeoutMin < 0 {
		return fmt.Errorf("job.default_timeout_min(%d) must not be minus value.", Job.DefaultTimeoutMin)
	}
	if Job.ConnectionTimeoutSec <= 0 {
		return fmt.Errorf("job.connection_timeout_sec(%d) must not be 0 or less.", Job.ConnectionTimeoutSec)
	}
	if Job.TimeTrackingSpanMin < 0 {
		return fmt.Errorf("job.time_tracking_span_min(%d) must not be minus value.", Job.TimeTrackingSpanMin)
	}
	if Job.AttemptLimit <= 0 {
		return fmt.Errorf("job.attempt_limit(%d) must not be 0 or less.", Job.AttemptLimit)
	}
	if Log.MaxSizeKB <= 0 {
		return fmt.Errorf("log.max_size_kb(%d) must not be 0 or less.", Log.MaxSizeKB)
	}
	if Log.MaxGeneration <= 0 {
		return fmt.Errorf("log.max_generation(%d) must not be 0 or less.", Log.MaxGeneration)
	}

	return nil
}
