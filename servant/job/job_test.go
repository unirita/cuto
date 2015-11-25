package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/unirita/cuto/message"
	"github.com/unirita/cuto/servant/config"
	"github.com/unirita/cuto/testutil"
)

var conf *config.ServantConfig
var testJobPath string

// ジョブログなどの掃除
func init() {
	time.Local = time.FixedZone("JST", 9*60*60)

	testJobPath = filepath.Join(testutil.GetBaseDir(), "servant", "job", "_testdata")
	err := os.Chdir(testJobPath)
	config.RootPath = testJobPath
	if err != nil {
		panic(err.Error())
	}
	configPath := filepath.Join(testJobPath, testServantIni)
	conf = config.ReadConfig(configPath)
	os.RemoveAll(conf.Dir.JoblogDir)
	err = os.Mkdir(conf.Dir.JoblogDir, 0755)
	if err != nil {
		panic(err.Error())
	}
}

// メッセージ・DB向けの時刻フォーマットをジョブログファイル名向けのフォーマットに変換する。
func stToLocalTimestamp(st string) string {
	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", st, time.UTC)
	if err != nil {
		panic("Unexpected time format: " + err.Error())
	}
	return t.Local().Format("20060102150405.000")
}

// ジョブログファイル名をフルパスで作成する。
// ”開始日(YYYYMMDD)\インスタンスID.ジョブ名（拡張子なし）.開始日時（yyyyMMddHHmmss.sss).log
func createJoblogFileName(req *message.Request, st string, nID int, jID string) string {
	job := filepath.Base(req.Path)
	if extpos := strings.LastIndex(job, "."); extpos != -1 {
		job = job[:extpos]
	}
	joblogFile := fmt.Sprintf("%v.%v.%v.%v.log", nID, job, jID, st)
	return filepath.Join(conf.Dir.JoblogDir, st[:8], joblogFile)
}
