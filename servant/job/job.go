// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package job

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/unirita/cuto/console"
	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/log"
	"github.com/unirita/cuto/message"
	"github.com/unirita/cuto/servant/config"
	"github.com/unirita/cuto/utctime"
)

// 実行ジョブ情報
type jobInstance struct {
	config          *config.ServantConfig // サーバントの設定情報
	nID             int                   // ネットワークID
	path            string                // ジョブファイル
	param           string                // 実行時パラメータ
	env             string                // 環境変数
	workDir         string                // 作業フォルダ
	wrnRC           int                   // 警告終了の戻り値
	wrnPtn          string                // 警告終了の文字列パターン
	errRC           int                   // 異常終了の戻り値
	errPtn          string                // 異常終了の文字列パターン
	timeout         int                   // 実行タイムアウトまでの時間（秒）
	jID             string                // ジョブID
	rc              int                   // ジョブの戻り値
	stat            int                   // ジョブステータス
	detail          string                // 異常終了時のメッセージ
	variable        string                // 変数情報
	st              string                // ジョブ開始日時
	et              string                // ジョブ終了日時
	joblog          string                // ジョブログ内容
	joblogFile      string                // ジョブログファイル名
	joblogTimestamp string                // ジョブログファイル名に使用するタイムスタンプ文字列
}

var (
	detailWarnRC  = "JOB-RC exceeded MAX-WarnRC."
	detailErrRC   = "JOB-RC exceeded MAX-ErrRC."
	detailWarnPtn = "JOB-OUTPUT matched Warning Pattern."
	detailErrPtn  = "JOB-OUTPUT matched Error Pattern."
)

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
	if err := job.do(stCh); err != nil {
		console.DisplayError("CTS019E", err)
		job.stat = db.ABNORMAL
		job.detail = err.Error()
		return job.createResponse()
	}

	console.Display("CTS011I", job.path, job.nID, job.jID, job.stat, job.rc)
	job.setVariableValue()
	return job.createResponse()
}

func (j *jobInstance) do(stCh chan<- string) error {
	isDockerJob := j.path == message.DockerTag
	cmd := j.createShell()
	if isDockerJob && cmd.Path == "" {
		return errors.New("Cannot execute job on Docker, because docker_command_path is lacked in servant.ini")
	}

	if err := j.run(cmd, stCh); err != nil {
		return err
	}

	if j.config.Job.DisuseJoblog == 0 {
		if err := j.writeJoblog(); err != nil {
			return err
		}
	}

	// RCからの結果と、出力MSGの結果を比較し、大きい方（異常の方）を採用する
	rcSt, rcMsg := j.judgeRC()
	ptnSt, ptnMsg := j.judgeJoblog()
	if rcSt > ptnSt {
		j.stat = rcSt
		j.detail = rcMsg
	} else {
		j.stat = ptnSt
		j.detail = ptnMsg
	}

	return nil
}

// ジョブファイルの拡張子を確認して、実行シェルを作成します。
func (j *jobInstance) createShell() *exec.Cmd {
	shell, params := j.organizePathAndParam()
	cmd := exec.Command(shell, params...)

	// 環境変数指定がない場合は、既存の物のみを追加する。
	if len(j.env) > 0 {
		envs := strings.Split(j.env, "+")
		cmd.Env = append(envs, os.Environ()...)
	} else {
		cmd.Env = os.Environ()
	}
	if len(j.workDir) > 0 {
		cmd.Dir = j.workDir
	} else {
		cmd.Dir = j.config.Dir.JobDir
	}

	return cmd
}

func (j *jobInstance) organizePathAndParam() (string, []string) {
	var shell string
	var params []string
	if j.path == message.DockerTag {
		shell = j.config.Job.DockerCommandPath
		params = paramSplit(j.param)
		// コンテナ上での実行ファイル名を用いてジョブログが作成されるよう、j.pathを上書き
		for index, param := range params {
			if param == "exec" {
				index += 2
				if index < len(params) {
					j.path = params[index]
				} else {
					j.path = ""
				}
				break
			}
		}
	} else {
		// ジョブファイル名のみの場合は、デフォルト場所を指定
		if !filepath.IsAbs(j.path) {
			j.path = filepath.Join(j.config.Dir.JobDir, j.path)
		}
		var paramStr string
		switch filepath.Ext(j.path) {
		case ".vbs":
			fallthrough
		case ".js":
			shell = "cscript"
			paramStr = fmt.Sprintf("/nologo %s %s", shellFormat(j.path), j.param)
		case ".jar":
			shell = "java"
			paramStr = fmt.Sprintf("-jar %s %s", shellFormat(j.path), j.param)
		case ".ps1":
			shell = "powershell"
			if sep := strings.IndexRune(j.path, ' '); sep != -1 {
				paramStr = fmt.Sprintf("\"& '%s' %s\"", j.path, j.param)
			} else {
				paramStr = fmt.Sprintf("%s %s", j.path, j.param)
			}
		default:
			shell = j.path
			paramStr = j.param
		}
		params = paramSplit(paramStr)
	}

	return shell, params
}

// ジョブ実行を行い、そのリターンコードを返す。
func (j *jobInstance) run(cmd *exec.Cmd, stCh chan<- string) error {
	isJoblogDisabled := j.config.Job.DisuseJoblog != 0
	outputBuffer := new(bytes.Buffer)

	if isJoblogDisabled {
		outputWriter := io.MultiWriter(os.Stdout, outputBuffer)
		cmd.Stdout = outputWriter
		cmd.Stderr = outputWriter
	} else {
		cmd.Stdout = outputBuffer
		cmd.Stderr = outputBuffer
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	startTime := utctime.Now()
	j.st = startTime.String() // ジョブ開始日時の取得
	j.joblogTimestamp = startTime.FormatLocaltime(utctime.NoDelimiter)
	stCh <- j.st

	console.Display("CTS010I", j.path, j.nID, j.jID, cmd.Process.Pid)

	err := j.waitCmdTimeout(cmd)
	j.et = utctime.Now().String() // ジョブ終了日時の取得

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
	j.joblog = outputBuffer.String()
	return err
}

func (j *jobInstance) waitCmdTimeout(cmd *exec.Cmd) error {
	if j.timeout == 0 {
		// timeoutが0の場合はタイムアウトなしでジョブ終了を待つ
		return cmd.Wait()
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)
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
func (j *jobInstance) judgeRC() (int, string) {
	if j.errRC > 0 {
		if j.errRC <= j.rc {
			return db.ABNORMAL, detailErrRC
		}
	}
	if j.wrnRC > 0 {
		if j.wrnRC <= j.rc {
			return db.WARN, detailWarnRC
		}
	}
	return db.NORMAL, ""
}

// ジョブログ結果を確認し、ステータスを返す。
func (j *jobInstance) judgeJoblog() (int, string) {
	if len(j.errPtn) > 0 {
		if strings.Contains(j.joblog, j.errPtn) {
			return db.ABNORMAL, detailErrPtn
		}
	}
	if len(j.wrnPtn) > 0 {
		if strings.Contains(j.joblog, j.wrnPtn) {
			return db.WARN, detailWarnPtn
		}
	}
	return db.NORMAL, ""
}

// ジョブログ結果を確認し、ステータスを返す。
// joblog内に指定された文字列が存在する場合は、それぞれのステータスを返します。
func (j *jobInstance) writeJoblog() error {
	// ジョブログファイル名の作成
	j.joblogFile = j.createJoblogFileName()
	log.Debug("joblogFile = ", j.joblogFile)

	// ファイルは存在しない場合の新規作成モード。
	file, err := os.OpenFile(j.joblogFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(j.joblog)
	return err
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
	joblogDir := filepath.Join(j.config.Dir.JoblogDir, j.joblogTimestamp[:8])
	if _, err := os.Stat(joblogDir); err != nil {
		os.Mkdir(joblogDir, 0777)
	}
	log.Debug("joblogDir = ", joblogDir)
	joblogFileName := fmt.Sprintf("%v.%v.%v.%v.log", j.nID, job, j.jID, j.joblogTimestamp)
	return filepath.Join(joblogDir, joblogFileName)
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
	res.JoblogFile = filepath.Base(j.joblogFile)
	return &res
}

// ジョブログファイルから変数情報を取得する。
func (j *jobInstance) setVariableValue() {
	reader := strings.NewReader(j.joblog)
	scanner := bufio.NewScanner(reader)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	j.variable = line
}
