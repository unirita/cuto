package config

import (
	"strings"
	"testing"
)

func generateTestConfig() *ServantConfig {
	c := new(ServantConfig)
	c.Sys.BindAddress = `0.0.0.0`
	c.Sys.BindPort = 0
	c.Job.HeartbeatSpanSec = 1
	c.Job.MultiProc = 1
	c.Job.DockerCommandPath = `/usr/bin/docker`
	c.Job.DisuseJoblog = 1
	c.Dir.JobDir = `.\jobscript`
	c.Dir.JoblogDir = `.\joblog`
	c.Dir.LogDir = `.\log`
	c.Log.OutputLevel = `info`
	c.Log.MaxSizeKB = 1
	c.Log.MaxGeneration = 1
	return c
}

func TestReadConfig_設定ファイルが開けない場合はデフォルト値をセットする(t *testing.T) {
	RootPath = `C:\cuto`
	FilePath = `noexistsfilepath`
	ReadConfig("")

	if Servant.Sys.BindAddress != defaultBindAddress {
		t.Errorf("bind_addressの設定値[%s]が想定と違っている。", Servant.Sys.BindAddress)
	}
	if Servant.Sys.BindPort != defaultBindPort {
		t.Errorf("bind_portの設定値[%d]が想定と違っている。", Servant.Sys.BindPort)
	}
	if Servant.Job.HeartbeatSpanSec != defaultHeartbeatSpanSec {
		t.Errorf("heartbeat_span_secの設定値[%d]が想定と違っている。", Servant.Job.HeartbeatSpanSec)
	}
	if Servant.Job.MultiProc != defaultMultiProc {
		t.Errorf("multi_procの設定値[%d]が想定と違っている。", Servant.Job.MultiProc)
	}
	if Servant.Job.DockerCommandPath != defaultDockerCommandPath {
		t.Errorf("docker_command_pathの設定値[%s]が想定と違っている。", Servant.Job.DockerCommandPath)
	}
	if Servant.Job.DisuseJoblog != defaultDisuseJoblog {
		t.Errorf("disuse_joblog[%d]が想定と違っている。", Servant.Job.DisuseJoblog)
	}
	if !strings.HasSuffix(Servant.Dir.JobDir, defaultJobDir) {
		t.Errorf("job_dirの設定値[%s]が想定と違っている。", Servant.Dir.JobDir)
	}
	if !strings.HasSuffix(Servant.Dir.JoblogDir, defaultJoblogDir) {
		t.Errorf("joblog_dirの設定値[%s]が想定と違っている。", Servant.Dir.JoblogDir)
	}
	if !strings.HasSuffix(Servant.Dir.LogDir, defaultLogDir) {
		t.Errorf("log_dirの設定値[%s]が想定と違っている。", Servant.Dir.LogDir)
	}
	if Servant.Log.OutputLevel != defaultOutputLevel {
		t.Errorf("output_levelの設定値[%s]が想定と違っている。", Servant.Log.OutputLevel)
	}
	if Servant.Log.MaxSizeKB != defaultMaxSizeKB {
		t.Errorf("max_size_kbの設定値[%d]が想定と違っている。", Servant.Log.MaxSizeKB)
	}
	if Servant.Log.MaxGeneration != defaultMaxGeneration {
		t.Errorf("max_generationの設定値[%d]が想定と違っている。", Servant.Log.MaxGeneration)
	}
}

func TestLoadReader_Readerから設定値を取得できる(t *testing.T) {
	conf := `
[sys]
bind_address='0.0.0.0'
bind_port=2015

[job]
multi_proc=20
heartbeat_span_sec=30
docker_command_path='/usr/bin/docker'
disuse_joblog=1

[dir]
joblog_dir='.\joblog'
job_dir='.\jobscript'
log_dir='.\log'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	cfg, err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	if cfg.Sys.BindAddress != `0.0.0.0` {
		t.Errorf("bind_addressの値[%s]が想定と違っている。", cfg.Sys.BindAddress)
	}
	if cfg.Sys.BindPort != 2015 {
		t.Errorf("bind_portの値[%d]が想定と違っている。", cfg.Sys.BindPort)
	}
	if cfg.Job.MultiProc != 20 {
		t.Errorf("multi_procの値[%d]が想定と違っている。", cfg.Job.MultiProc)
	}
	if cfg.Job.HeartbeatSpanSec != 30 {
		t.Errorf("heartbeat_span_sec[%d]が想定と違っている。", cfg.Job.HeartbeatSpanSec)
	}
	if cfg.Job.DockerCommandPath != `/usr/bin/docker` {
		t.Errorf("docker_command_pathの設定値[%s]が想定と違っている。", cfg.Job.DockerCommandPath)
	}
	if cfg.Job.DisuseJoblog != 1 {
		t.Errorf("disuse_joblog[%d]が想定と違っている。", cfg.Job.DisuseJoblog)
	}
	if cfg.Dir.JoblogDir != `.\joblog` {
		t.Errorf("joblog_dirの値[%s]が想定と違っている。", cfg.Dir.JoblogDir)
	}
	if cfg.Dir.JobDir != `.\jobscript` {
		t.Errorf("job_dirの値[%s]が想定と違っている。", cfg.Dir.JobDir)
	}
	if cfg.Dir.LogDir != `.\log` {
		t.Errorf("log_dirの値[%s]が想定と違っている。", cfg.Dir.LogDir)
	}
	if cfg.Log.OutputLevel != `info` {
		t.Errorf("output_levelの値[%s]が想定と違っている。", cfg.Log.OutputLevel)
	}
	if cfg.Log.MaxSizeKB != 10240 {
		t.Errorf("max_size_kb[%d]は想定と違っている。", cfg.Log.MaxSizeKB)
	}
	if cfg.Log.MaxGeneration != 2 {
		t.Errorf("max_generation[%d]は想定と違っている。", cfg.Log.MaxGeneration)
	}
}

func TestLoadReader_CUTOROOTタグを展開できる(t *testing.T) {
	conf := `
[sys]
bind_address='0.0.0.0'
bind_port=2015

[job]
multi_proc=20
heartbeat_span_sec=30
docker_command_path='/usr/bin/docker'
disuse_joblog=1

[dir]
joblog_dir='<CUTOROOT>\joblog'
job_dir='<CUTOROOT>\jobscript'
log_dir='<CUTOROOT>\log'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	cfg, err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	if cfg.Dir.JoblogDir == `<CUTOROOT>\joblog` {
		t.Errorf("joblog_dir内のCUTOROOTタグが展開されていない。")
	}
	if cfg.Dir.JobDir == `<CUTOROOT>\jobscript` {
		t.Errorf("job_dir内のCUTOROOTタグが展開されていない。")
	}
	if cfg.Dir.LogDir == `<CUTOROOT>\log` {
		t.Errorf("log_dir内のCUTOROOTタグが展開されていない。")
	}
}

func TestLoadReader_tomlの書式に沿っていない場合はエラーが発生する(t *testing.T) {
	conf := `
[sys]
bind_address=0.0.0.0
bind_port=2015

[job]
multi_proc=20
heartbeat_span_sec=30
docker_command_path='/usr/bin/docker'
disuse_joblog=1

[dir]
joblog_dir='.\joblog'
job_dir='.\jobscript'
log_dir='.\log'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	_, err := loadReader(r)
	if err == nil {
		t.Error("エラーが発生しなかった")
	}
}

func TestDetectError_設定値が正常な場合はエラーが発生しない(t *testing.T) {
	c := generateTestConfig()
	if err := c.DetectError(); err != nil {
		t.Errorf("想定以外のエラーが発生した: %s", err)
	}
}

func TestDetectError_設定値が正常な場合はエラーが発生しない_ポート番号最大値(t *testing.T) {
	c := generateTestConfig()
	c.Sys.BindPort = 65535
	if err := c.DetectError(); err != nil {
		t.Errorf("想定以外のエラーが発生した: %s", err)
	}
}

func TestDetectError_デフォルトのポート番号が負の値の場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Sys.BindPort = -1
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_デフォルトのポート番号が65535を超える場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Sys.BindPort = 65536
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ハートビートメッセージ送信間隔が0以下の場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Job.HeartbeatSpanSec = 0
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ジョブ多重度が0以下の場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Job.MultiProc = 0
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログファイル最大サイズが0以下の場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Log.MaxSizeKB = 0
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログファイル最大世代数が0以下の場合はエラー(t *testing.T) {
	c := generateTestConfig()
	c.Log.MaxGeneration = 0
	if err := c.DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}
