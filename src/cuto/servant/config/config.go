// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"

	"cuto/console"
	"cuto/util"
)

const (
	defaultBindAddress      = `0.0.0.0`
	defaultBindPort         = 2015
	defaultHeartbeatSpanSec = 30
	defaultMultiProc        = 20
	defaultJobDir           = `jobscript`
	defaultJoblogDir        = `joblog`
	defaultLogDir           = `log`
	defaultOutputLevel      = `info`
	defaultMaxSizeKB        = 10240
	defaultMaxGeneration    = 2
)

const dirName = "bin"
const fileName = "servant.ini"

// 設定情報のデフォルト値を設定する。
func DefaultServantConfig() *ServantConfig {
	cfg := new(ServantConfig)
	cfg.Sys.BindAddress = defaultBindAddress
	cfg.Sys.BindPort = defaultBindPort
	cfg.Job.HeartbeatSpanSec = defaultHeartbeatSpanSec
	cfg.Job.MultiProc = defaultMultiProc
	cfg.Dir.JobDir = defaultJobDir
	cfg.Dir.JoblogDir = defaultJoblogDir
	cfg.Dir.LogDir = defaultLogDir
	cfg.Log.OutputLevel = defaultOutputLevel
	cfg.Log.MaxSizeKB = defaultMaxSizeKB
	cfg.Log.MaxGeneration = defaultMaxGeneration

	return cfg

}

// サーバント設定情報
type ServantConfig struct {
	Sys sysSection
	Job jobSection
	Dir dirSection
	Log logSection
}

// サーバント設定のsysセクション
type sysSection struct {
	BindAddress string `toml:"bind_address"`
	BindPort    int    `toml:"bind_port"`
}

// サーバント設定のjobセクション
type jobSection struct {
	MultiProc        int `toml:"multi_proc"`
	HeartbeatSpanSec int `toml:"heartbeat_span_sec"`
}

// サーバント設定のdirセクション
type dirSection struct {
	JoblogDir string `toml:"joblog_dir"`
	JobDir    string `toml:"job_dir"`
	LogDir    string `toml:"log_dir"`
}

// 設定ファイルのlogセクション
type logSection struct {
	OutputLevel   string `toml:"output_level"`
	MaxSizeKB     int    `toml:"max_size_kb"`
	MaxGeneration int    `toml:"max_generation"`
}

var Servant *ServantConfig
var FilePath string
var RootPath string

func init() {
	RootPath = util.GetRootPath()
	//	FilePath = fmt.Sprintf("%s%c%s%c%s", RootPath, os.PathSeparator, dirName, os.PathSeparator, fileName)
	FilePath = fmt.Sprintf(".%c%s%c%s", os.PathSeparator, dirName, os.PathSeparator, fileName)
}

// 設定ファイルを読み込む
// 読み込みに失敗する場合はDefaultServantConfig関数でデフォルト値を設定する。
//
// 戻り値: 設定値を格納したServantConfig構造体オブジェクト
func ReadConfig(configPath string) *ServantConfig {
	var err error
	if len(configPath) > 0 {
		FilePath = configPath
	}
	Servant, err = loadFile(FilePath)
	if err != nil {
		console.Display("CTS004W", FilePath)
		Servant = DefaultServantConfig()
	}

	Servant.convertFullpath()
	return Servant
}

// 設定をリロードする。
//
// 戻り値: 設定値を格納したServantConfig構造体オブジェクト
func ReloadConfig() *ServantConfig {
	//@TODO 今はとりあえずファイル読むだけ
	return ReadConfig(FilePath)
}

// 設定値のエラー検出を行う。
//
// 戻り値: エラー情報
func (c *ServantConfig) DetectError() error {
	if c.Sys.BindPort < 0 || 65535 < c.Sys.BindPort {
		return fmt.Errorf("sys.bind_port(%d) must be within the range 0 and 65535.", c.Sys.BindPort)
	}
	if c.Job.HeartbeatSpanSec <= 0 {
		return fmt.Errorf("job.heartbeat_span_sec(%d) must not be 0 or less.", c.Job.HeartbeatSpanSec)
	}
	if c.Job.MultiProc <= 0 {
		return fmt.Errorf("job.multi_proc(%d) must not be 0 or less.", c.Job.MultiProc)
	}
	if c.Log.MaxSizeKB <= 0 {
		return fmt.Errorf("log.max_size_kb(%d) must not be 0 or less.", c.Log.MaxSizeKB)
	}
	if c.Log.MaxGeneration <= 0 {
		return fmt.Errorf("log.max_generation(%d) must not be 0 or less.", c.Log.MaxGeneration)
	}

	return nil
}

func loadFile(filePath string) (*ServantConfig, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return loadReader(f)
}

func loadReader(reader io.Reader) (*ServantConfig, error) {
	sc := new(ServantConfig)
	if _, err := toml.DecodeReader(reader, sc); err != nil {
		return nil, err
	}

	return sc, nil
}

// 設定値の相対パスを絶対パスへ変換する。
func (s *ServantConfig) convertFullpath() {
	if strings.HasPrefix(s.Dir.JoblogDir, ".\\") || strings.HasPrefix(s.Dir.JoblogDir, "./") {
		s.Dir.JoblogDir = strings.Replace(s.Dir.JoblogDir, ".", RootPath, 1)
	}
	if strings.HasPrefix(s.Dir.JobDir, ".\\") || strings.HasPrefix(s.Dir.JobDir, "./") {
		s.Dir.JobDir = strings.Replace(s.Dir.JobDir, ".", RootPath, 1)
	}
	if strings.HasPrefix(s.Dir.LogDir, ".\\") || strings.HasPrefix(s.Dir.LogDir, "./") {
		s.Dir.LogDir = strings.Replace(s.Dir.LogDir, ".", RootPath, 1)
	}
}
