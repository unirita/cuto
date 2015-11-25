package job

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func createTestJobInstance() *jobInstance {
	j := new(jobInstance)
	j.config = config.DefaultServantConfig()
	j.nID = 1234
	j.jID = "JID"
	j.path = "test.sh"
	j.param = "param1 param2"
	j.env = "env1 env2"
	j.workDir = ""
	j.wrnRC = 5
	j.wrnPtn = "warn"
	j.errRC = 10
	j.errPtn = "error"
	j.timeout = 100
	return j
}

func TestJobCreateShell_Normal(t *testing.T) {
	j := createTestJobInstance()
	cmd := j.createShell()
	expectedPath := filepath.Join(j.config.Dir.JobDir, "test.sh")
	if cmd.Path != expectedPath {
		t.Errorf("cmd.Path => %s, wants %s", cmd.Path, expectedPath)
	}
	if len(cmd.Args) != 3 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 3)
	}
	if cmd.Args[1] != "param1" {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], "param1")
	}
	if cmd.Args[2] != "param2" {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], "param2")
	}
}

func TestJobCreateShell_Docker_CommandPathIsSet(t *testing.T) {
	j := createTestJobInstance()
	j.config.Job.DockerCommandPath = "/usr/bin/docker"
	j.path = message.DockerTag
	j.param = "exec containerName command param1 param2"
	cmd := j.createShell()
	if cmd.Path != "/usr/bin/docker" {
		t.Errorf("cmd.Path => %s, wants %s", cmd.Path, "/usr/bin/docker")
	}
	if len(cmd.Args) != 6 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 6)
	}
	if cmd.Args[1] != "exec" {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], "exec")
	}
	if cmd.Args[2] != "containerName" {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], "containerName")
	}
	if cmd.Args[3] != "command" {
		t.Errorf("cmd.Args[3] => %s, wants %s", cmd.Args[3], "command")
	}
	if cmd.Args[4] != "param1" {
		t.Errorf("cmd.Args[4] => %s, wants %s", cmd.Args[4], "param1")
	}
	if cmd.Args[5] != "param2" {
		t.Errorf("cmd.Args[5] => %s, wants %s", cmd.Args[5], "param2")
	}
	if j.path != "command" {
		t.Errorf("j.path => %s, wants %s", j.path, "command")
	}
}

func TestJobCreateShell_Docker_CommandPathIsNotSet(t *testing.T) {
	j := createTestJobInstance()
	j.config.Job.DockerCommandPath = ""
	j.path = message.DockerTag
	cmd := j.createShell()
	if cmd.Path != "" {
		t.Errorf("cmd.Path => %s, wants %s", cmd.Path, "")
	}
}

func TestJobCreateShell_VBScript(t *testing.T) {
	j := createTestJobInstance()
	j.path = "test.vbs"
	cmd := j.createShell()
	if !strings.Contains(cmd.Path, "cscript") {
		t.Errorf("cmd.Path must contains '%s', but it did not.", "cscript")
		t.Logf("cmd.Path => %s", cmd.Path)
	}
	if len(cmd.Args) != 5 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 5)
	}
	if cmd.Args[1] != "/nologo" {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], "/nologo")
	}
	expectedPath := filepath.Join(j.config.Dir.JobDir, "test.vbs")
	if cmd.Args[2] != expectedPath {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], expectedPath)
	}
	if cmd.Args[3] != "param1" {
		t.Errorf("cmd.Args[3] => %s, wants %s", cmd.Args[3], "param1")
	}
	if cmd.Args[4] != "param2" {
		t.Errorf("cmd.Args[4] => %s, wants %s", cmd.Args[4], "param2")
	}
}

func TestJobCreateShell_JScript(t *testing.T) {
	j := createTestJobInstance()
	j.path = "test.js"
	cmd := j.createShell()
	if !strings.Contains(cmd.Path, "cscript") {
		t.Errorf("cmd.Path must contains '%s', but it did not.", "cscript")
		t.Logf("cmd.Path => %s", cmd.Path)
	}
	if len(cmd.Args) != 5 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 5)
	}
	if cmd.Args[1] != "/nologo" {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], "/nologo")
	}
	expectedPath := filepath.Join(j.config.Dir.JobDir, "test.js")
	if cmd.Args[2] != expectedPath {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], expectedPath)
	}
	if cmd.Args[3] != "param1" {
		t.Errorf("cmd.Args[3] => %s, wants %s", cmd.Args[3], "param1")
	}
	if cmd.Args[4] != "param2" {
		t.Errorf("cmd.Args[4] => %s, wants %s", cmd.Args[4], "param2")
	}
}

func TestJobCreateShell_JAR(t *testing.T) {
	j := createTestJobInstance()
	j.path = "test.jar"
	cmd := j.createShell()
	if !strings.Contains(cmd.Path, "java") {
		t.Errorf("cmd.Path must contains '%s', but it did not.", "java")
		t.Logf("cmd.Path => %s", cmd.Path)
	}
	if len(cmd.Args) != 5 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 5)
	}
	if cmd.Args[1] != "-jar" {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], "-jar")
	}
	expectedPath := filepath.Join(j.config.Dir.JobDir, "test.jar")
	if cmd.Args[2] != expectedPath {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], expectedPath)
	}
	if cmd.Args[3] != "param1" {
		t.Errorf("cmd.Args[3] => %s, wants %s", cmd.Args[3], "param1")
	}
	if cmd.Args[4] != "param2" {
		t.Errorf("cmd.Args[4] => %s, wants %s", cmd.Args[4], "param2")
	}
}

func TestJobCreateShell_PowerShell(t *testing.T) {
	j := createTestJobInstance()
	j.path = "test.ps1"
	cmd := j.createShell()
	if !strings.Contains(cmd.Path, "powershell") {
		t.Errorf("cmd.Path must contains '%s', but it did not.", "powershell")
		t.Logf("cmd.Path => %s", cmd.Path)
	}
	if len(cmd.Args) != 4 {
		t.Fatalf("len(cmd.Args) => %d, wants %d", len(cmd.Args), 5)
	}
	expectedPath := filepath.Join(j.config.Dir.JobDir, "test.ps1")
	if cmd.Args[1] != expectedPath {
		t.Errorf("cmd.Args[1] => %s, wants %s", cmd.Args[1], expectedPath)
	}
	if cmd.Args[2] != "param1" {
		t.Errorf("cmd.Args[2] => %s, wants %s", cmd.Args[2], "param1")
	}
	if cmd.Args[3] != "param2" {
		t.Errorf("cmd.Args[3] => %s, wants %s", cmd.Args[3], "param2")
	}
}
