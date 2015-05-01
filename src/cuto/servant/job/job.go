// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package job

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"os/exec"

	"cuto/console"
	"cuto/db"
	"cuto/log"
	"cuto/message"
	"cuto/servant/config"
	"cuto/util"
)

// 実行ジョブ情報
type jobInstance struct {
	config     *config.ServantConfig // サーバントの設定情報
	nID        int                   // ネットワークID
	path       string                // ジョブファイル
	param      string                // 実行時パラメータ
	env        string                // 環境変数
	workDir    string                // 作業フォルダ
	wrnRC      int                   // 警告終了の戻り値
	wrnPtn     string                // 警告終了の文字列パターン
	errRC      int                   // 異常終了の戻り値
	errPtn     string                // 異常終了の文字列パターン
	timeout    int                   // 実行タイムアウトまでの時間（秒）
	jID        string                // ジョブID
	rc         int                   // ジョブの戻り値
	stat       int                   // ジョブステータス
	detail     string                // 異常終了時のメッセージ
	variable   string                // 変数情報
	st         string                // ジョブ開始日時
	et         string                // ジョブ終了日時
	joblog     string                // ジョブログ内容
	joblogFile string                // ジョブログファイル名
}

// 実行ジョブ情報のコンストラクタ
func newJobInstance(req *message.Request, conf *config.ServantConfig) *jobInstance {
	job := new(jobInstance)
	job.config = conf
	job.nID = req.NID
	job.path = req.Path
	job.param = req.Param
	job.env = req.Env
	job.workDir = req.Workspace
	job.wrnRC = req.WarnRC
	job.wrnPtn = req.WarnStr
	job.errRC = req.ErrRC
	job.errPtn = req.ErrStr
	job.timeout = req.Timeout
	job.jID = req.JID

	return job
}

// ジョブの実行要求を受け付けて実行する。
//
// param : req マスタからの要求メッセージ。
//
// param : conf サーバントの設定情報。
//
// param : stCh スタート時刻送信用チャンネル
//
// return : マスタへ返信するメッセージ。
func DoJobRequest(req *message.Request, conf *config.ServantConfig, stCh chan<- string) *message.Response {
	job := newJobInstance(req, conf)

	cmd := job.createShell()

	err := job.run(cmd, stCh)
	if err != nil {
		log.Error(err)
		job.stat = db.ABNORMAL
		job.detail = err.Error()
	} else {
		rcSt := job.judgeRC()
		ptnSt, err := job.writeFileAndJodgeJoblog()
		if err != nil {
			log.Error(err)
			job.stat = db.ABNORMAL
			job.detail = err.Error()
		} else {
			if rcSt > ptnSt {
				job.stat = rcSt
			} else {
				job.stat = ptnSt
			}
		}
		console.Display("CTS011I", job.path, job.nID, job.jID, job.stat, job.rc)
		if job.stat != db.ABNORMAL {
			job.setVariableValue()
		}
	}

	return job.createResponse()
}

// ジョブファイルの拡張子を確認して、実行シェルを作成します。
func (j *jobInstance) createShell() *exec.Cmd {
	var shell, param, script string
	// ジョブファイル名のみの場合は、デフォルト場所を指定
	if !strings.Contains(j.path, "\\") && !strings.Contains(j.path, "/") {
		script = fmt.Sprintf("%s%c%s", j.config.Dir.JobDir, os.PathSeparator, j.path)
	} else {
		script = j.path
	}
	// 拡張子に応じた、実行シェルを作成
	if strings.HasSuffix(j.path, ".vbs") || strings.HasSuffix(j.path, ".js") { // WSH
		shell = "cscript"
		param = fmt.Sprintf("/nologo %s %s", script, j.param)
	} else if strings.HasSuffix(j.path, ".jar") { // JAVA
		shell = "java"
		param = fmt.Sprintf("-jar %s %s", script, j.param)
	} else if strings.HasSuffix(j.path, ".ps1") { // PowerShell
		shell = "powershell"
		param = fmt.Sprintf("%s %s", script, j.param)
	} else { // bat or exe
		shell = script
		param = j.param
	}
	params := strings.Split(param, " ")

	cmd := exec.Command(shell, params...)
	envs := strings.Split(j.env, "+")
	cmd.Env = append(os.Environ(), envs...)
	if len(j.workDir) > 0 {
		cmd.Dir = j.workDir
	} else {
		cmd.Dir = j.config.Dir.JobDir
	}

	return cmd
}

// ジョブ実行を行い、そのリターンコードを返す。
func (j *jobInstance) run(cmd *exec.Cmd, stCh chan<- string) error {
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if err := cmd.Start(); err != nil {
		return err
	}
	j.st = util.DateJoblogFormat(time.Now()) // ジョブ開始日時の取得
	stCh <- j.st

	console.Display("CTS010I", j.path, j.nID, j.jID, cmd.Process.Pid)

	err := j.waitCmdTimeout(cmd)
	j.et = util.DateJoblogFormat(time.Now()) // ジョブ終了日時の取得

	if err != nil {
		if e2, ok := err.(*exec.ExitError); ok {
			if s, ok := e2.Sys().(syscall.WaitStatus); ok {
				j.rc = s.ExitStatus()
				err = nil
			} else {
				j.detail = errors.New("Unimplemented for system where exec.ExitError.Sys() is not syscall.WaitStatus.").Error()
			}
		}
	} else {
		j.rc = 0
	}
	j.joblog = stdout.String()
	return err
}

func (j *jobInstance) waitCmdTimeout(cmd *exec.Cmd) error {
	if j.timeout == 0 {
		// timeoutが0の場合はタイムアウトなしでジョブ終了を待つ
		return cmd.Wait()
	}

	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Wait()
	}()

	t := time.Duration(j.timeout) * time.Second
	select {
	case err := <-ch:
		return err
	case <-time.After(t):
		cmd.Process.Kill()
		return errors.New("Process timeout.")
	}

	return nil
}

// ジョブのRCを確認し、statを返す。
// ジョブのRCが指定されたRC以上の場合は、それぞれのステータスを返します。
func (j *jobInstance) judgeRC() int {
	if j.errRC > 0 {
		if j.errRC <= j.rc {
			return db.ABNORMAL
		}
	}
	if j.wrnRC > 0 {
		if j.wrnRC <= j.rc {
			return db.WARN
		}
	}
	return db.NORMAL
}

// ジョブログ結果を確認し、ステータスを返す。
// joblog内に指定された文字列が存在する場合は、それぞれのステータスを返します。
func (j *jobInstance) writeFileAndJodgeJoblog() (int, error) {
	// ジョブログファイル名の作成
	j.joblogFile = j.createJoblogFileName()
	log.Debug("joblogFile = ", j.joblogFile)

	// ファイルは存在しない場合の新規作成モード。
	file, err := os.OpenFile(j.joblogFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	_, err = file.WriteString(j.joblog)
	if err != nil {
		return 0, err
	}

	if len(j.errPtn) > 0 {
		if strings.Contains(j.joblog, j.errPtn) {
			return db.ABNORMAL, nil
		}
	}
	if len(j.wrnPtn) > 0 {
		if strings.Contains(j.joblog, j.wrnPtn) {
			return db.WARN, nil
		}
	}
	return db.NORMAL, nil
}

// ジョブログファイル名をフルパスで作成する。
// ”開始日(YYYYMMDD)\インスタンスID.ジョブ名（拡張子なし）.開始日時（yyyyMMddHHmmss.sss).log
func (j *jobInstance) createJoblogFileName() string {
	var job string // ジョブ名（拡張子なし）の取得
	if strings.LastIndex(j.path, "\\") != -1 {
		tokens := strings.Split(j.path, "\\")
		job = tokens[len(tokens)-1]
	} else if strings.LastIndex(j.path, "/") != -1 {
		tokens := strings.Split(j.path, "/")
		job = tokens[len(tokens)-1]
	} else {
		job = j.path
	}
	if extpos := strings.LastIndex(job, "."); extpos != -1 {
		job = job[:extpos]
	}
	// 開始日フォルダの作成
	joblogDir := fmt.Sprintf("%v%c%v", j.config.Dir.JoblogDir, os.PathSeparator, j.st[:8])
	if _, err := os.Stat(joblogDir); err != nil {
		os.Mkdir(joblogDir, 0666)
	}
	log.Debug("joblogDir = ", joblogDir)
	return fmt.Sprintf("%v%c%v.%v.%v.log", joblogDir, os.PathSeparator, j.nID, job, j.st)
}

// レスポンスメッセージの作成
func (j *jobInstance) createResponse() *message.Response {
	var res message.Response
	res.NID = j.nID
	res.JID = j.jID
	res.RC = j.rc
	res.Stat = j.stat
	res.Detail = j.detail
	res.Var = j.variable
	res.St = j.st
	res.Et = j.et
	return &res
}

// ジョブログファイルから変数情報を取得する。
func (j *jobInstance) setVariableValue() {
	file, err := os.Open(j.joblogFile)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	j.variable = line
}
