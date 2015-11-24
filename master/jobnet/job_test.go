package jobnet

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/db/tx"
	"github.com/unirita/cuto/master/config"
	"github.com/unirita/cuto/message"
)

func newTestNetwork() *Network {
	n, _ := NewNetwork("test")
	n.Result = &tx.ResultMap{JobnetResult: nil, Jobresults: make(tx.JobMap)}

	dbpath := getTestDBPath()
	conn, err := db.Open(dbpath)
	if err != nil {
		panic(err)
	}
	n.Result.SetConnection(conn)
	return n
}

func testSendRequest_Normal(host string, port int, reqMsg string, stCh chan<- string) (string, error) {
	req := new(message.Request)
	req.ParseJSON("reqMsg")

	res := new(message.Response)
	res.NID = req.NID
	res.JID = req.JID
	res.RC = 0
	res.Stat = 1
	res.Detail = ""
	res.Var = "testvar"
	res.St = "2015-04-01 12:34:56.789"
	res.Et = "2015-04-01 12:35:46.123"

	resMsg, _ := res.GenerateJSON()
	return resMsg, nil
}

func testSendRequest_Abnormal(host string, port int, reqMsg string, stCh chan<- string) (string, error) {
	req := new(message.Request)
	req.ParseJSON("reqMsg")

	res := new(message.Response)
	res.NID = req.NID
	res.JID = req.JID
	res.RC = 1
	res.Stat = 9
	res.Detail = "testerror"
	res.Var = "testvar"
	res.St = "2015-04-01 12:34:56.789"
	res.Et = "2015-04-01 12:35:46.123"

	resMsg, _ := res.GenerateJSON()
	return resMsg, nil
}

func testSendRequest_Error(host string, port int, reqMsg string, stCh chan<- string) (string, error) {
	return "", fmt.Errorf("senderror")
}

func testSendRequest_ErrorAfterSt(host string, port int, reqMsg string, stCh chan<- string) (string, error) {
	stCh <- "2015-04-01 12:34:56.789"
	time.Sleep(time.Millisecond * 50)
	return "", fmt.Errorf("senderror")
}

func testSendRequest_NotJSON(host string, port int, reqMsg string, stCh chan<- string) (string, error) {
	return "notjson", nil
}

func TestNewJob_ジョブ構造体にデフォルト値がセットされる(t *testing.T) {
	id := "1234"
	name := "testjob"

	job, err := NewJob(id, name, nil)

	if job.id != id {
		t.Errorf("IDとして%sが期待されるのに対し、%sがセットされました。", id, job.id)
	}
	if job.Name != name {
		t.Errorf("ジョブ名として%dが期待されるのに対し、%dがセットされました。", name, job.Name)
	}
	if err != nil {
		t.Errorf("ジョブ名[%v]に禁止文字が使用されていませんが、エラーを返しました。。", name)
	}
}

func TestNewJob_ジョブ名に禁止文字が含まれるとエラー(t *testing.T) {
	id := "1234"
	name := "te:stjob"

	_, err := NewJob(id, name, nil)

	if err == nil {
		t.Errorf("ジョブ名[%v]に禁止文字が組まれているにかかわらず、エラーを返しませんでした。", name)
	}
}

func TestJobType_ジョブ構造体のノードタイプを取得できる(t *testing.T) {
	id := "1234"
	name := "testjob"
	job, _ := NewJob(id, name, nil)

	if job.Type() != ELM_JOB {
		t.Errorf("間違ったノードタイプ[%v]が取得された。", job.Type())
	}
}

func TestJobAddNext_ジョブ構造体に後続エレメントを追加できる(t *testing.T) {
	j1 := new(Job)
	j2 := new(Job)
	g1 := new(Gateway)

	if err := j1.AddNext(j2); err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if j1.Next != j2 {
		t.Error("後続ジョブの追加に失敗した。")
	}
	if err := j2.AddNext(g1); err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if j2.Next != g1 {
		t.Error("後続ゲートウェイの追加に失敗した。")
	}
}

func TestJobAddNext_後続エレメントを複数追加しようとした場合はエラー(t *testing.T) {
	j1 := new(Job)
	j2 := new(Job)
	j3 := new(Job)
	j1.AddNext(j2)

	if err := j1.AddNext(j3); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestJobHasNext_ジョブ構造体の後続エレメントの有無をチェックできる(t *testing.T) {
	j1 := new(Job)
	j2 := new(Job)
	g1 := new(Gateway)

	if j1.HasNext() {
		t.Error("後続エレメントが無いのにも関わらず、HasNextがtrueを返した")
	}

	j1.Next = j2
	if !j1.HasNext() {
		t.Error("後続ジョブがあるのにも関わらず、HasNextがfalseを返した")
	}

	j2.Next = g1
	if !j2.HasNext() {
		t.Error("後続ゲートウェイがあるのにも関わらず、HasNextがfalseを返した")
	}
}

func TestJobExecute_レスポンスにエラーが無いケース(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Next = j2

	j1.sendRequest = testSendRequest_Normal
	next, err := j1.Execute()
	if err != nil {
		t.Fatalf("想定外のエラーが発生: %s", err)
	}
	if next != j2 {
		t.Errorf("次に実行されるのとは違うノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.StartDate != "2015-04-01 12:34:56.789" {
		t.Errorf("ジョブ実行結果のStartDate[%s]は想定と違っている。", jobres.StartDate)
	}
	if jobres.EndDate != "2015-04-01 12:35:46.123" {
		t.Errorf("ジョブ実行結果のEndDate[%s]は想定と違っている。", jobres.EndDate)
	}
	if jobres.Status != 1 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail != "" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Rc != 0 {
		t.Errorf("ジョブ実行結果のRc[%d]は想定と違っている。", jobres.Rc)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
	if jobres.Variable != "testvar" {
		t.Errorf("ジョブ実行結果のVariable[%s]は想定と違っている。", jobres.Variable)
	}
}

func TestJobExecute_使用できない変数を使用したケース(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Param = "$MJNOEXISTS_RC$"
	j1.Next = j2

	j1.sendRequest = testSendRequest_Normal
	next, err := j1.Execute()
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("nilが返される想定に対し、ノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.Status != 9 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail == "" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
}

func TestJobExecute_ジョブが異常終了したケース(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Next = j2

	j1.sendRequest = testSendRequest_Abnormal
	next, err := j1.Execute()
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("nilが返される想定に対し、ノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.StartDate != "2015-04-01 12:34:56.789" {
		t.Errorf("ジョブ実行結果のStartDate[%s]は想定と違っている。", jobres.StartDate)
	}
	if jobres.EndDate != "2015-04-01 12:35:46.123" {
		t.Errorf("ジョブ実行結果のEndDate[%s]は想定と違っている。", jobres.EndDate)
	}
	if jobres.Status != 9 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail != "testerror" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Rc != 1 {
		t.Errorf("ジョブ実行結果のRc[%d]は想定と違っている。", jobres.Rc)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
	if jobres.Variable != "testvar" {
		t.Errorf("ジョブ実行結果のVariable[%s]は想定と違っている。", jobres.Variable)
	}
}

func TestJobExecute_リクエスト送信に失敗したケース(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Next = j2

	j1.sendRequest = testSendRequest_Error
	next, err := j1.Execute()
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("nilが返される想定に対し、ノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.Status != 9 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail != "senderror" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
}

func TestJobExecute_リクエスト送信に失敗したケース_失敗前にスタート時刻を受け取った場合(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Next = j2

	j1.sendRequest = testSendRequest_ErrorAfterSt
	next, err := j1.Execute()
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("nilが返される想定に対し、ノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.StartDate != "2015-04-01 12:34:56.789" {
		t.Errorf("ジョブ実行結果のStartDate[%s]は想定と違っている。", jobres.StartDate)
	}
	if jobres.Status != 9 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail != "senderror" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
}

func TestJobExecute_レスポンスがJSON形式でないケース(t *testing.T) {
	config.Job.AttemptLimit = 1
	n := newTestNetwork()
	j1, _ := NewJob("jobid1", "job1", n)
	j2, _ := NewJob("jobid1", "job1", n)
	j1.Node = "testnode"
	j1.Port = 1234
	j1.Next = j2

	j1.sendRequest = testSendRequest_NotJSON
	next, err := j1.Execute()
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
	if next != nil {
		t.Errorf("nilが返される想定に対し、ノード[%s]が返された。", next.ID())
	}

	jobres, ok := n.Result.Jobresults[j1.id]
	if !ok {
		t.Fatal("ジョブ実行結果がセットされなかった。")
	}
	if jobres.ID != n.ID {
		t.Errorf("ジョブ実行結果のID[%d]は想定と違っている。", jobres.ID)
	}
	if jobres.JobId != j1.id {
		t.Errorf("ジョブ実行結果のJobId[%s]は想定と違っている。", jobres.JobId)
	}
	if jobres.JobName != j1.Name {
		t.Errorf("ジョブ実行結果のJobName[%s]は想定と違っている。", jobres.JobName)
	}
	if jobres.Status != 9 {
		t.Errorf("ジョブ実行結果のStatus[%d]は想定と違っている。", jobres.Status)
	}
	if jobres.Detail == "" {
		t.Errorf("ジョブ実行結果のDetail[%s]は想定と違っている。", jobres.Detail)
	}
	if jobres.Node != "testnode" {
		t.Errorf("ジョブ実行結果のNode[%s]は想定と違っている。", jobres.Node)
	}
	if jobres.Port != 1234 {
		t.Errorf("ジョブ実行結果のPort[%d]は想定と違っている。", jobres.Port)
	}
}

func TestCreateJoblogFileName(t *testing.T) {
	n := newTestNetwork()
	j, _ := NewJob("jobid1", "job1", n)
	j.FilePath = `C:\test\testjob.bat`

	res := new(message.Response)
	res.NID = n.ID
	res.JID = j.ID()
	res.RC = 1
	res.Stat = 9
	res.Detail = "testerror"
	res.Var = "testvar"
	res.St = "2015-04-01 12:34:56.789"
	res.Et = "2015-04-01 12:35:46.123"

	actual := j.createJoblogFileName(res)
	expected := strconv.Itoa(n.ID) + `.testjob.jobid1.20150401123456.789.log`
	if actual != expected {
		t.Errorf("ジョブログファイル名[%s]は想定値[%s]と違っている。", actual, expected)
	}
}

func TestExplodeNodeString_NoContainer(t *testing.T) {
	node := "testhost"
	host, containerHost, containerName := explodeNodeString(node)
	if host != "testhost" {
		t.Errorf("host => %s, wants %s", host, "testhost")
	}
	if containerHost != "" {
		t.Error("containerHost must be empty, but it was not.")
		t.Log("containerHost:", containerHost)
	}
	if containerName != "" {
		t.Error("containerName must be empty, but it was not.")
		t.Log("containerName:", containerName)
	}
}

func TestExplodeNodeString_ContainerWithoutHost(t *testing.T) {
	node := "testhost>category/name"
	host, containerHost, containerName := explodeNodeString(node)
	if host != "testhost" {
		t.Errorf("host => %s, wants %s", host, "testhost")
	}
	if containerHost != "" {
		t.Error("containerHost must be empty, but it was not.")
		t.Log("containerHost:", containerHost)
	}
	if containerName != "category/name" {
		t.Errorf("containerName => %s, wants %s", containerName, "category/name")
	}
}

func TestExplodeNodeString_ContainerWithHost(t *testing.T) {
	node := "testhost>tcp://hostname/category/name"
	host, containerHost, containerName := explodeNodeString(node)
	if host != "testhost" {
		t.Errorf("host => %s, wants %s", host, "testhost")
	}
	if containerHost != "tcp://hostname" {
		t.Errorf("containerHost => %s, wants %s", containerHost, "tcp://hostname")
	}
	if containerName != "category/name" {
		t.Errorf("containerName => %s, wants %s", containerName, "category/name")
	}
}
