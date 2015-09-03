// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package jobnet

import (
	"fmt"
	"strings"
	"time"

	"cuto/console"
	"cuto/db"
	"cuto/db/tx"
	"cuto/log"
	"cuto/master/config"
	"cuto/master/remote"
	"cuto/message"
	"cuto/util"
)

type sendFunc func(string, int, string, chan<- string) (string, error)

// ジョブを表す構造体
type Job struct {
	id            string   // ジョブID
	Name          string   // ジョブ名
	Node          string   // ノード
	Port          int      // ポート番号
	FilePath      string   // ジョブファイル
	Param         string   // ジョブ引き渡しパラメータ
	Env           string   // ジョブ実行に渡す環境変数
	Workspace     string   // ジョブ実行時の作業フォルダ
	WrnRC         int      // 警告終了と判断する戻り値の下限値
	WrnPtn        string   // 警告終了と判断するジョブの出力メッセージ
	ErrRC         int      // 異常終了と判断する戻り値の下限値
	ErrPtn        string   // 異常終了と判断するジョブの出力メッセージ
	Timeout       int      // ジョブ実行時間のタイムアウト
	SecondaryNode string   // セカンダリサーバントのノード
	SecondaryPort int      // セカンダリサーバントのポート番号
	Next          Element  // 次ノード
	Instance      *Network // ネットワーク情報構造体のポインタ
	sendRequest   sendFunc // リクエスト送信メソッド
	IsRerunJob    bool     // リランジョブであるかどうか
}

// Job構造体のコンストラクタ関数。
//
// param ; id  ジョブネットワークID。
//
// param : name  ジョブネットワーク名。
//
// param : nwk  ジョブネットワーク構造体。
//
// return : ジョブ情報構造体。
//
// return : エラー情報。
func NewJob(id string, name string, nwk *Network) (*Job, error) {
	if util.JobnameHasInvalidRune(name) {
		return nil, fmt.Errorf("Job name[%s] includes forbidden character.", name)
	}
	job := new(Job)
	job.id = id
	job.Name = name
	job.Instance = nwk
	job.sendRequest = remote.SendRequest
	return job, nil
}

// IDを取得する
//
// return : ジョブID.
func (j *Job) ID() string {
	return j.id
}

// ノードタイプを取得する
//
// return : ノードタイプ
func (j *Job) Type() elementType {
	return ELM_JOB
}

// 後続エレメントの追加を行う。
//
// param : 追加する要素情報。
func (j *Job) AddNext(e Element) error {
	if j.Next != nil {
		return fmt.Errorf("ServiceTask cannot connect with over 1 element.")
	}
	j.Next = e

	return nil
}

// 後続エレメントの有無を調べる。
//
// return : 要素が存在する場合はtrueを返す。
func (j *Job) HasNext() bool {
	return j.Next != nil
}

// 拡張ジョブ情報のデフォルト値をセットする
func (j *Job) SetDefaultEx() {
	if j.Node == "" {
		j.Node = config.Job.DefaultNode
	}
	if j.Port == 0 {
		j.Port = config.Job.DefaultPort
	}
	if j.FilePath == "" {
		j.FilePath = j.Name
	}
	if j.Timeout < 0 {
		j.Timeout = config.Job.DefaultTimeoutMin * 60
	}
}

// ジョブ実行リクエストをservantへ送信する。
//
// return : 次の実行ノード
//
// return : エラー情報。
func (j *Job) Execute() (Element, error) {
	if j.IsRerunJob {
		jobres := j.Instance.Result.Jobresults[j.id]
		if jobres.Status == db.NORMAL || jobres.Status == db.WARN {
			j.resumeJobValue()
			return j.Next, nil
		} else {
			result, err := j.requestLatestJobResult()
			if err != nil {
				return nil, j.abnormalEnd(err)
			}

			switch result.Stat {
			case db.RUNNING:
				return nil, fmt.Errorf("Job ID [%s] still running.", j.id)
			case db.NORMAL:
				fallthrough
			case db.WARN:
				j.updateNormalEndResult(result)
				return j.Next, nil
			default:
				j.changeStatusRunning()
			}
		}
	}
	res, err := j.executeRequest()
	if err != nil {
		return nil, j.abnormalEnd(err)
	}
	defer j.end(res)

	if isAbnormalEnd(res) {
		console.Display("CTM026I", j.createJoblogFileName(res), j.Node)
		return nil, fmt.Errorf("")
	}

	return j.Next, nil
}

func (j *Job) executeRequest() (*message.Response, error) {
	req := new(message.Request)
	req.NID = j.Instance.ID
	req.JID = j.ID()
	req.Path = j.FilePath
	req.Param = j.Param
	req.Env = j.Env
	req.Workspace = j.Workspace
	req.WarnRC = j.WrnRC
	req.WarnStr = j.WrnPtn
	req.ErrRC = j.ErrRC
	req.ErrStr = j.ErrPtn
	req.Timeout = j.Timeout

	if !j.IsRerunJob {
		j.start(req)
	}

	err := req.ExpandMasterVars()
	if err != nil {
		return nil, err
	}

	reqMsg, err := req.GenerateJSON()
	if err != nil {
		return nil, err
	}

	stCh := make(chan string, 1)
	go j.waitAndSetResultStartDate(stCh)

	timerEndCh := make(chan struct{}, 1)
	go j.startTimer(timerEndCh)
	defer close(timerEndCh)

	resMsg, err := j.sendRequestWithRetry(reqMsg, stCh)
	if j.isNecessaryToRetry(err) && j.SecondaryNode != "" {
		j.useSecondaryNode()
		console.Display("CTM028W", j.Name, j.SecondaryNode)
		resMsg, err = j.sendSecondaryRequestWithRetry(reqMsg, stCh)
	}

	close(stCh)
	if err != nil {
		return nil, err
	}

	res := new(message.Response)
	err = res.ParseJSON(resMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (j *Job) requestLatestJobResult() (*message.JobResult, error) {
	chk := new(message.JobCheck)
	chk.NID = j.Instance.ID
	chk.JID = j.ID()

	chkMsg, err := chk.GenerateJSON()
	if err != nil {
		return nil, err
	}

	resultMsg, err := j.sendResultCheckRequestWithRetry(chkMsg)
	if err != nil {
		return nil, err
	}

	result := new(message.JobResult)
	err = result.ParseJSON(resultMsg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ジョブ実行リクエストを送信する。
// 送信失敗時には必要な回数だけリトライを行う。
func (j *Job) sendRequestWithRetry(reqMsg string, stCh chan<- string) (string, error) {
	limit := config.Job.AttemptLimit
	var resMsg string
	var err error
	for i := 0; i < limit; i++ {
		if i != 0 {
			console.Display("CTM027W", j.Name, i, limit-1)
		}

		resMsg, err = j.sendRequest(j.Node, j.Port, reqMsg, stCh)
		if !j.isNecessaryToRetry(err) {
			break
		}
	}

	return resMsg, err
}

// セカンダリサーバントへのジョブ実行リクエストを送信する。
// 送信失敗時には必要な回数だけリトライを行う。
func (j *Job) sendSecondaryRequestWithRetry(reqMsg string, stCh chan<- string) (string, error) {
	limit := config.Job.AttemptLimit
	var resMsg string
	var err error
	for i := 0; i < limit; i++ {
		if i != 0 {
			console.Display("CTM027W", j.Name, i, limit-1)
		}

		resMsg, err = j.sendRequest(j.SecondaryNode, j.SecondaryPort, reqMsg, stCh)
		if !j.isNecessaryToRetry(err) {
			break
		}
	}

	return resMsg, err
}

func (j *Job) sendResultCheckRequestWithRetry(chkMsg string) (string, error) {
	stCh := make(chan string, 1)
	defer close(stCh)

	limit := config.Job.AttemptLimit
	var resultMsg string
	var err error
	for i := 0; i < limit; i++ {
		if i != 0 {
			console.Display("CTM027W", j.Name, i, limit-1)
		}

		resultMsg, err = j.sendRequest(j.Node, j.Port, chkMsg, stCh)
		if err == nil {
			break
		}
	}

	return resultMsg, err
}

// リトライの必要があるかを判定する。
// 判定条件：リクエスト送信でエラーが発生しており、スタート時刻がセットされていないこと
func (j *Job) isNecessaryToRetry(err error) bool {
	if err != nil {
		jobres, exist := j.Instance.Result.Jobresults[j.id]
		if !exist {
			return true
		} else if jobres.StartDate == "" {
			return true
		}
	}

	return false
}

// responseメッセージrのステータスを参照し、ジョブが異常終了している場合はtrueを返す。
// それ以外はfalseを返す。
func isAbnormalEnd(r *message.Response) bool {
	if r.Stat == db.ABNORMAL {
		return true
	}
	return false
}

// ジョブの開始処理を行う。
func (j *Job) start(req *message.Request) {
	jobres := db.NewJobResult(int(j.Instance.ID))
	jobres.JobId = j.ID()
	jobres.JobName = j.Name
	jobres.Node = j.Node
	jobres.Port = j.Port
	jobres.Status = db.RUNNING

	j.Instance.Result.Jobresults[j.ID()] = jobres
	tx.InsertJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)

	console.Display("CTM023I", j.Name, j.Node, j.Instance.ID, j.id)
}

// 実行ノードをセカンダリノードに変更する。
func (j *Job) useSecondaryNode() {
	jobres, exist := j.Instance.Result.Jobresults[j.id]
	if !exist {
		log.Error(fmt.Errorf("Job result[id = %s] is unregisted.", j.id))
		return
	}

	jobres.Node = j.SecondaryNode
	jobres.Port = j.SecondaryPort
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)
}

// ジョブ実行結果にジョブの開始時刻をセットする。
func (j *Job) waitAndSetResultStartDate(stCh <-chan string) {
	st := <-stCh
	if len(st) == 0 {
		// （主にチャネルがクローズされることにより）空文字列が送られてきた場合は何もしない。
		return
	}
	log.Debug(fmt.Sprintf("JOB[%s] StartDate[%s]", j.Name, st))

	jobres, exist := j.Instance.Result.Jobresults[j.id]
	if !exist {
		log.Error(fmt.Errorf("Job result[id = %s] is unregisted.", j.id))
		return
	}
	jobres.StartDate = st
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)
}

// ジョブの終了メッセージから、ジョブ状態の更新を行う。
func (j *Job) end(res *message.Response) {
	var jobres *db.JobResult
	var exist bool

	if jobres, exist = j.Instance.Result.Jobresults[j.id]; !exist {
		log.Error(fmt.Errorf("Job result[id = %s] is unregisted.", j.id))
		return
	}
	jobres.StartDate = res.St
	jobres.EndDate = res.Et
	jobres.Status = res.Stat
	jobres.Rc = res.RC
	jobres.Detail = res.Detail
	jobres.Variable = res.Var

	message.AddJobValue(j.Name, res)
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)

	var st string
	switch jobres.Status {
	case db.NORMAL:
		st = db.ST_NORMAL
	case db.WARN:
		st = db.ST_WARN
	default:
		st = db.ST_ABNORMAL
	}
	if jobres.Status != db.ABNORMAL {
		console.Display("CTM024I", j.Name, j.Node, j.Instance.ID, j.id, st)
	} else {
		console.Display("CTM025W", j.Name, j.Node, j.Instance.ID, j.id, st, jobres.Detail)
	}
}

// サーバントへ送受信失敗した場合の異常終了処理
func (j *Job) abnormalEnd(err error) error {
	jobres, exist := j.Instance.Result.Jobresults[j.id]
	if !exist {
		return fmt.Errorf("Job result[id = %s] is unregisted.", j.id)
	}
	jobres.Status = db.ABNORMAL
	jobres.Detail = err.Error()
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)

	console.Display("CTM025W", j.Name, j.Node, j.Instance.ID, j.id, jobres.Status, jobres.Detail)
	return err
}

func (j *Job) resumeJobValue() {
	jobres := j.Instance.Result.Jobresults[j.id]

	res := new(message.Response)
	res.JID = j.id
	res.RC = jobres.Rc
	res.St = jobres.StartDate
	res.Et = jobres.EndDate
	res.Var = jobres.Variable
	message.AddJobValue(j.Name, res)
}

func (j *Job) updateNormalEndResult(result *message.JobResult) {
	jobres, exists := j.Instance.Result.Jobresults[j.id]
	if !exists {
		log.Error(fmt.Errorf("Job result[id = %s] is unregisted.", j.id))
		return
	}

	jobres.Status = result.Stat
	jobres.Rc = result.RC
	jobres.StartDate = result.St
	jobres.EndDate = result.Et
	jobres.Detail = ""
	jobres.Variable = result.Var
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)

	j.resumeJobValue()
}

func (j *Job) changeStatusRunning() {
	jobres, exists := j.Instance.Result.Jobresults[j.id]
	if !exists {
		log.Error(fmt.Errorf("Job result[id = %s] is unregisted.", j.id))
		return
	}

	jobres.Status = db.RUNNING
	tx.UpdateJob(j.Instance.Result.GetConnection(), jobres, &j.Instance.localMutex)
}

func (j *Job) startTimer(endCh chan struct{}) {
	span := config.Job.TimeTrackingSpanMin
	if span == 0 {
		// 出力間隔の設定が0の場合は出力しない。
		return
	}

	rapTime := 0
	for {
		select {
		case <-time.After(time.Duration(span) * time.Minute):
			rapTime += span
			console.Display("CTM022I", j.Name, rapTime)
		case <-endCh:
			return
		}
	}
}

func (j *Job) createJoblogFileName(r *message.Response) string {
	// ジョブ名（拡張子なし）の取得
	job := j.FilePath
	if strings.LastIndex(job, "\\") != -1 {
		tokens := strings.Split(job, "\\")
		job = tokens[len(tokens)-1]
	} else if strings.LastIndex(job, "/") != -1 {
		tokens := strings.Split(job, "/")
		job = tokens[len(tokens)-1]
	}
	if extpos := strings.LastIndex(job, "."); extpos != -1 {
		job = job[:extpos]
	}

	timestamp := r.St
	timestamp = strings.Replace(timestamp, "-", "", -1)
	timestamp = strings.Replace(timestamp, " ", "", -1)
	timestamp = strings.Replace(timestamp, ":", "", -1)

	return fmt.Sprintf("%v.%v.%v.%v.log", j.Instance.ID, job, j.ID(), timestamp)
}
