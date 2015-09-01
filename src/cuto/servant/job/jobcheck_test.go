package job

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cuto/message"
	"cuto/servant/config"
	"cuto/utctime"
)

func init() {
	time.Local = time.FixedZone("JST", 9*60*60)
}

func getJobCheckTestConfig() *config.ServantConfig {
	testDir := filepath.Join(os.Getenv("GOPATH"), "test", "cuto", "servant", "job", "jobcheck")

	conf := new(config.ServantConfig)
	conf.Dir.JoblogDir = filepath.Join(testDir, "joblog")
	conf.Dir.LogDir = filepath.Join(testDir, "log")
	return conf
}

func TestDoJobResultCheck_Base(t *testing.T) {
	chk := &message.JobCheck{
		Type:    "jobcheck",
		Version: "1.2.3",
		NID:     1,
		JID:     "job1",
	}

	result := DoJobResultCheck(chk, getJobCheckTestConfig())
	if result == nil {
		t.Fatalf("DoJobResultCheck() returned nil.")
	}
	if result.NID != 1 {
		t.Errorf("result.NID => %d, wants %d", result.NID, 1)
	}
	if result.JID != "job1" {
		t.Errorf("result.JID => %s, wants %s", result.JID, "job1")
	}
	if result.RC != 5 {
		t.Errorf("result.RC => %d, wants %d", result.RC, 5)
	}
	if result.Stat != 2 {
		t.Errorf("result.Stat => %d, wants %d", result.Stat, 5)
	}
	if result.Var != "testvar" {
		t.Errorf("result.Var => %s, wants %s", result.Var, "testvar")
	}
	if result.St != "2015-08-01 03:05:25.123" {
		t.Errorf("result.St => %s, wants %s", result.St, "2015-08-01 03:05:25.123")
	}
	if result.Et != "2015-08-01 03:34:56.789" {
		t.Errorf("result.Et => %s, wants %s", result.Et, "2015-08-01 03:34:56.789")
	}
}

func TestDoJobResultCheck_RegardNoRecordJobAsExecuting(t *testing.T) {
	chk := &message.JobCheck{
		Type:    "jobcheck",
		Version: "1.2.3",
		NID:     1,
		JID:     "noexists",
	}

	result := DoJobResultCheck(chk, getJobCheckTestConfig())
	if result == nil {
		t.Fatalf("DoJobResultCheck() returned nil.")
	}
	if result.NID != 1 {
		t.Errorf("result.NID => %d, wants %d", result.NID, 1)
	}
	if result.JID != "noexists" {
		t.Errorf("result.JID => %s, wants %s", result.JID, "noexists")
	}
	if result.Stat != 0 {
		t.Errorf("result.Stat => %d, wants %d", result.Stat, 0)
	}
}

func TestSearchEndRecordFromLog(t *testing.T) {
	conf := getJobCheckTestConfig()
	logPath := filepath.Join(conf.Dir.LogDir, "servant.log")
	record, err := searchJobEndRecordFromLog(logPath, 1, "job1")
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}

	expected := `2015-08-01 12:34:56.789 [16594] [INF] CTS011I JOB [/home/cuto/testjob] ENDED. INSTANCE [1] ID [job1] STATUS [2] RC [5].`
	if record != expected {
		t.Errorf("Record is not expected value.")
		t.Log("Actual:")
		t.Log(record)
		t.Log("Expected: ")
		t.Log(expected)
	}
}

func TestSearchLatestJoblog(t *testing.T) {
	conf := getJobCheckTestConfig()
	et, _ := utctime.Parse(utctime.Default, "2015-08-01 03:34:56.789")
	path, err := searchLatestJoblog(conf.Dir.JoblogDir, 1, "job1", et)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}

	if !strings.HasSuffix(path, "1.testjob.job1.20150801120525.123.log") {
		t.Error("Unexpected path was got:")
		t.Log(path)
	}
}
