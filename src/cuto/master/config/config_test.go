package config

import (
	"strings"
	"testing"
)

func generateTestConfig() {
	Job.DefaultNode = `localhost`
	Job.DefaultPort = 0
	Job.DefaultTimeoutMin = 0
	Job.ConnectionTimeoutSec = 1
	Job.TimeTrackingSpanMin = 10
	Job.AttemptLimit = 5
	Dir.JobnetDir = `.\jobnet`
	Dir.LogDir = `.\log`
	DB.DBFile = `.\data\cuto.sqlite`
	Log.OutputLevel = `info`
	Log.MaxSizeKB = 1
	Log.MaxGeneration = 1
}

func TestLoad_存在しないファイルをロードしようとした場合はエラー(t *testing.T) {
	if err := Load("noexistfilepath"); err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestLoadByReader_Readerから設定値を取得できる(t *testing.T) {
	conf := `
[job]
default_node='localhost'
default_port=2015
default_timeout_min=30
connection_timeout_sec=60
time_tracking_span_min=10
attempt_limit=5

[dir]
jobnet_dir='jobnet'
log_dir='log'

[db]
db_file='cute.db'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	if Job.DefaultNode != `localhost` {
		t.Errorf("default_nodeの値[%s]は想定と違っている。", Job.DefaultNode)
	}
	if Job.DefaultPort != 2015 {
		t.Errorf("default_portの値[%d]は想定と違っている。", Job.DefaultPort)
	}
	if Job.DefaultTimeoutMin != 30 {
		t.Errorf("default_timeout_minの値[%d]は想定と違っている。", Job.DefaultTimeoutMin)
	}
	if Job.ConnectionTimeoutSec != 60 {
		t.Errorf("connection_timeout_secの値[%d]は想定と違っている。", Job.ConnectionTimeoutSec)
	}
	if Job.TimeTrackingSpanMin != 10 {
		t.Errorf("time_tracking_span_minの値[%d]は想定と違っている。", Job.TimeTrackingSpanMin)
	}
	if Job.AttemptLimit != 5 {
		t.Errorf("attempt_limitの値[%d]は想定と違っている。", Job.AttemptLimit)
	}
	if Dir.JobnetDir != `jobnet` {
		t.Errorf("jobnet_dirの値[%s]は想定と違っている。", Dir.JobnetDir)
	}
	if Dir.LogDir != `log` {
		t.Errorf("log_dirの値[%s]は想定と違っている。", Dir.LogDir)
	}
	if DB.DBFile != `cute.db` {
		t.Errorf("db_fileの値[%s]は想定と違っている。", DB.DBFile)
	}
	if Log.OutputLevel != `info` {
		t.Errorf("output_levelの値[%s]が想定と違っている。", Log.OutputLevel)
	}
	if Log.MaxSizeKB != 10240 {
		t.Errorf("max_size_kb[%d]は想定と違っている。", Log.MaxSizeKB)
	}
	if Log.MaxGeneration != 2 {
		t.Errorf("max_generation[%d]は想定と違っている。", Log.MaxGeneration)
	}
}

func TestLoadByReader_CUTOROOTタグを展開できる(t *testing.T) {
	conf := `
[job]
default_node='localhost'
default_port=2015
default_timeout_min=30
connection_timeout_sec=60
time_tracking_span_min=10
attempt_limit=5

[dir]
jobnet_dir='<CUTOROOT>/jobnet'
log_dir='<CUTOROOT>/log'

[db]
db_file='cute.db'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}
	if Dir.JobnetDir == `<CUTOROOT>/jobnet` {
		t.Errorf("jobnet_dir内の<CUTOROOT>が置換されていない")
	}
	if Dir.LogDir == `<CUTOROOT>/log` {
		t.Errorf("log_dir内の<CUTOROOT>が置換されていない")
	}
}

func TestLoadByReader_AttemptLimitが指定されない場合のデフォルト値(t *testing.T) {
	conf := `
[job]
default_node='localhost'
default_port=2015
default_timeout_min=30
connection_timeout_sec=60
time_tracking_span_min=10

[dir]
jobnet_dir='jobnet'
log_dir='log'

[db]
db_file='cute.db'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	if Job.AttemptLimit != 1 {
		t.Errorf("attempt_limitの値[%d]は想定と違っている。", Job.AttemptLimit)
	}
}

func TestLoadByReader_tomlの書式に沿っていない場合はエラーが発生する(t *testing.T) {
	conf := `
[job]
default_node=localhost
default_port=2015
default_timeout_min=30
connection_timeout_sec=60
time_tracking_span_min=10
attempt_limit=5

[dir]
jobnet_dir='jobnet'
log_dir='log'

[db]
db_file='cute.db'

[log]
output_level='info'
max_size_kb=10240
max_generation=2
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err == nil {
		t.Error("エラーが発生しなかった")
	}
}

func TestDetectError_設定内容にエラーが無い場合はnilを返す_ポート番号最小値(t *testing.T) {
	generateTestConfig()
	if err := DetectError(); err != nil {
		t.Errorf("想定外のエラーが発生した： %s", err)
	}
}

func TestDetectError_設定内容にエラーが無い場合はnilを返す_ポート番号最大値(t *testing.T) {
	generateTestConfig()
	Job.DefaultPort = 65535
	if err := DetectError(); err != nil {
		t.Errorf("想定外のエラーが発生した： %s", err)
	}
}

func TestDetectError_デフォルトのポート番号が負の値の場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.DefaultPort = -1
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_デフォルトのポート番号が65535を超える場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.DefaultPort = 65536
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_デフォルトの実行タイムアウト時間が負の値の場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.DefaultTimeoutMin = -1
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_接続タイムアウト時間が0以下の場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.ConnectionTimeoutSec = 0
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_経過時間表示間隔が負の値の場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.TimeTrackingSpanMin = -1
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_最大試行回数が0以下の場合はエラー(t *testing.T) {
	generateTestConfig()
	Job.AttemptLimit = 0
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログファイル最大サイズが0以下の場合はエラー(t *testing.T) {
	generateTestConfig()
	Log.MaxSizeKB = 0
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログファイル最大世代数が0以下の場合はエラー(t *testing.T) {
	generateTestConfig()
	Log.MaxGeneration = 0
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}
