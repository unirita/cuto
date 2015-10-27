package job

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/unirita/cuto/message"
	"github.com/unirita/cuto/servant/config"
	"github.com/unirita/cuto/utctime"
)

func DoJobResultCheck(chk *message.JobCheck, conf *config.ServantConfig) *message.JobResult {
	result := new(message.JobResult)
	result.NID = chk.NID
	result.JID = chk.JID

	logPath := filepath.Join(conf.Dir.LogDir, "servant.log")
	endRecord, err := searchJobEndRecordFromLog(logPath, result.NID, result.JID)
	if err != nil || len(endRecord) == 0 {
		return createUnexecutedResult(chk.NID, chk.JID)
	}

	status, err := extractStatusFromRecord(endRecord)
	if err != nil || status == 0 {
		return createErrorResult(chk.NID, chk.JID)
	}
	result.Stat = status

	et, err := extractTimestampFromRecord(endRecord)
	if err != nil {
		return createErrorResult(chk.NID, chk.JID)
	}
	result.Et = et.Format(utctime.Default)

	rc, err := extractRCFromRecord(endRecord)
	if err != nil {
		return createErrorResult(chk.NID, chk.JID)
	}
	result.RC = rc

	joblog, err := searchLatestJoblog(conf.Dir.JoblogDir, chk.NID, chk.JID, et)
	if err != nil {
		return createErrorResult(chk.NID, chk.JID)
	}

	st, err := extractTimestampFromJoblog(joblog)
	if err != nil {
		return createErrorResult(chk.NID, chk.JID)
	}
	result.St = st.Format(utctime.Default)

	variable, err := extractVariableFromJoblog(joblog)
	if err != nil {
		return createErrorResult(chk.NID, chk.JID)
	}
	result.Var = variable

	return result
}

func searchJobEndRecordFromLog(path string, nid int, jid string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	matchStr := fmt.Sprintf(
		`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} \[\d+\] \[INF\] CTS01[01]I.*INSTANCE \[%d\] ID \[%s\]`,
		nid, jid)
	matcher := regexp.MustCompile(matchStr)
	var endRecord string
	s := bufio.NewScanner(file)
	for s.Scan() {
		record := s.Text()
		if matcher.MatchString(record) {
			endRecord = record
		}
	}

	return endRecord, nil
}

func extractTimestampFromRecord(record string) (utctime.UTCTime, error) {
	timestampStr := record[:len(utctime.Default)]
	return utctime.ParseLocaltime(utctime.Default, timestampStr)
}

func extractStatusFromRecord(record string) (int, error) {
	if strings.Contains(record, "CTS010I") {
		return 0, nil
	}

	finder := regexp.MustCompile(`STATUS \[\d+\]`)
	statusStr := finder.FindString(record)
	if len(statusStr) < 9 {
		return 0, errors.New("Could not extract status.")
	}
	status, err := strconv.Atoi(statusStr[8 : len(statusStr)-1])
	if err != nil {
		return 0, errors.New("Could not extract status as int.")
	}

	return status, nil
}

func extractRCFromRecord(record string) (int, error) {
	finder := regexp.MustCompile(`RC \[\d+\]`)
	rcStr := finder.FindString(record)
	if len(rcStr) < 5 {
		return 0, errors.New("Could not extract RC.")
	}
	rc, err := strconv.Atoi(rcStr[4 : len(rcStr)-1])
	if err != nil {
		return 0, errors.New("Could not extract RC as int.")
	}

	return rc, nil
}

func searchLatestJoblog(joblogDir string, nid int, jid string, et utctime.UTCTime) (string, error) {
	dirNames := make([]string, 2)
	dirNames[0] = et.FormatLocaltime(utctime.Date8Num)
	dirNames[1] = et.AddDays(-1).FormatLocaltime(utctime.Date8Num)

	matchStr := fmt.Sprintf(`^%d\.[^.]+\.%s\.`, nid, jid)
	matcher := regexp.MustCompile(matchStr)
	for _, dirName := range dirNames {
		dir := filepath.Join(joblogDir, dirName)
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			continue
		}

		var path string
		for _, fileInfo := range fileInfos {
			fileName := fileInfo.Name()
			if matcher.MatchString(fileName) {
				path = fileName
			}
		}
		if path != "" {
			return filepath.Join(dir, path), nil
		}
	}

	return "", errors.New("Could not find joblog file.")
}

func extractTimestampFromJoblog(joblog string) (utctime.UTCTime, error) {
	parts := strings.Split(joblog, ".")
	timestampStr := strings.Join(parts[len(parts)-3:len(parts)-1], ".")
	return utctime.ParseLocaltime(utctime.NoDelimiter, timestampStr)
}

func extractVariableFromJoblog(joblog string) (string, error) {
	file, err := os.Open(joblog)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var variable string
	s := bufio.NewScanner(file)
	for s.Scan() {
		variable = s.Text()
	}

	return variable, nil
}

func createUnexecutedResult(nid int, jid string) *message.JobResult {
	result := new(message.JobResult)
	result.NID = nid
	result.JID = jid
	result.Stat = -1
	return result
}

func createErrorResult(nid int, jid string) *message.JobResult {
	result := new(message.JobResult)
	result.NID = nid
	result.JID = jid
	return result
}
