package jobnet

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/unirita/cuto/master/config"
	"github.com/unirita/cuto/master/jobnet/parser"
	"github.com/unirita/cuto/testutil"
)

// ※補足事項
// testJob構造体およびその生成関数generateTestJobはpath_test.goで定義されています。

type undefinedElement struct {
	next Element
}

func (u *undefinedElement) ID() string        { return "undefined" }
func (u *undefinedElement) Type() elementType { return ELM_JOB }
func (u *undefinedElement) AddNext(e Element) error {
	u.next = e
	return nil
}
func (u *undefinedElement) HasNext() bool {
	if u.next == nil {
		return false
	}
	return true
}
func (u *undefinedElement) Execute() (Element, error) {
	return u.next, nil
}

func getTestDBPath() string {
	return filepath.Join(testutil.GetBaseDir(),
		"master", "jobnet", "_testdata", "test.sqlite")
}

func loadTestConfig() {
	config.Job.DefaultNode = `localhost`
	config.Job.DefaultPort = 2015
	config.Job.DefaultTimeoutMin = 30
	config.Dir.JobnetDir = `jobnet`
	config.DB.DBFile = getTestDBPath()
}

func TestNewNetwork_各メンバの初期値をセットできる(t *testing.T) {
	loadTestConfig()

	n, _ := NewNetwork("test")

	if n.Name != "test" {
		t.Errorf("セットされたネットワーク名[%s]が間違っています。", n.Name)
	}
	if n.elements == nil {
		t.Error("elementsがmakeされていません。")
	}
	expectedNwkPath := fmt.Sprintf("jobnet%ctest.bpmn", os.PathSeparator)
	if n.MasterPath != expectedNwkPath {
		t.Errorf("セットされたネットワーク定義ファイルパス[%s]が間違っています。", n.MasterPath)
	}
	expectedExPath := fmt.Sprintf("jobnet%ctest.csv", os.PathSeparator)
	if n.JobExPath != expectedExPath {
		t.Errorf("セットされた拡張ジョブ定義ファイルパス[%s]が間違っています。", n.JobExPath)
	}
}

func TestLoadNetwork_定義ファイルが存在しない場合はnilを返す(t *testing.T) {
	n := LoadNetwork("noexistnetwork")
	if n != nil {
		t.Errorf("返された値がnilではない。")
	}
}

func TestLoadElements_要素をロードできる(t *testing.T) {
	bpmn := `
<definitions>
    <process>
	    <startEvent id="start"/>
	    <endEvent id="end"/>
		<serviceTask id="job" name="jobname"/>
		<sequenceFlow sourceRef="start" targetRef="job"/>
		<sequenceFlow sourceRef="job" targetRef="end"/>
	</process>
</definitions>`

	n, _ := NewNetwork("test")
	err := n.LoadElements(strings.NewReader(bpmn))
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if len(n.elements) == 0 {
		t.Fatal("Network.elementsが空のままになっている。")
	}
}

func TestLoadElements_パースに失敗した場合はエラー(t *testing.T) {
	bpmn := `CannotParseString`

	n, _ := NewNetwork("test")
	err := n.LoadElements(strings.NewReader(bpmn))
	if err == nil {
		t.Error("エラーが発生しなかった。")
	}
	if len(n.elements) != 0 {
		t.Error("Network.elementsは空のままであるはずが、値がセットされた。")
	}
}

func TestSetElements_JobとGatewayを追加できる(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 1),
		Gateway: make([]parser.ParallelGateway, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	elements := nwk.elements
	if len(elements) != 2 {
		t.Fatalf("要素数は2個であるはずだが、%d個存在する。", len(elements))
	}

	task1, ok := elements["task1"]
	if !ok {
		t.Fatal("task1がセットされなかった。")
	}
	if task1.ID() != "task1" {
		t.Errorf("ジョブtask1のidはtask1であるはずだが、%sがセットされている。", task1.ID())
	}

	gw1, ok := elements["gw1"]
	if !ok {
		t.Fatal("gw1がセットされなかった。")
	}
	if gw1.ID() != "gw1" {
		t.Errorf("ゲートウェイgw1のidはgw1であるはずだが、%sがセットされている。", gw1.ID())
	}
}

func TestSetElements_Jobの先行関係を設定できる(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 2),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Flow[0] = parser.SequenceFlow{From: "task1", To: "task2"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	task1 := nwk.elements["task1"].(*Job)
	if task1.Next == nil {
		t.Fatal("task1の後続ジョブが設定されていない。")
	}

	task2, ok := task1.Next.(*Job)
	if !ok {
		t.Fatal("task1の後続ジョブがJob型になっていない。")
	}
	if task2.ID() != "task2" {
		t.Error("task1の後続ジョブはtask2であるはずが、%sになっている。", task2.ID())
	}
}

func TestSetElements_Gatewayの先行関係を設定できる(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 3),
		Gateway: make([]parser.ParallelGateway, 1),
		Flow:    make([]parser.SequenceFlow, 3),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Flow[0] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[1] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task3"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	task1 := nwk.elements["task1"].(*Job)
	if task1.Next == nil {
		t.Fatal("task1の後続ゲートウェイが設定されていない。")
	}

	gw1, ok := task1.Next.(*Gateway)
	if !ok {
		t.Fatal("task1の後続ゲートウェイがGateway型になっていない。")
	}
	if gw1.ID() != "gw1" {
		t.Error("task1の後続ゲートウェイはgw1であるはずが、%sになっている。", gw1.ID())
	}

	if len(gw1.Nexts) != 2 {
		t.Fatal("gw1の後続ジョブ数が2つになるはずが、違っている。")
	}
	if gw1.Nexts[0] == nil || gw1.Nexts[1] == nil {
		t.Fatal("gw1の後続ジョブのに設定されるはずのないnilが設定されている。")
	}

	task2, ok := gw1.Nexts[0].(*Job)
	if !ok {
		t.Fatal("gw1の1つめの後続ジョブがJob型になっていない。")
	}
	if task2.ID() != "task2" {
		t.Error("gw1の1つめの後続ジョブはtask2であるはずが、%sになっている。", task2.ID())
	}

	task3, ok := gw1.Nexts[1].(*Job)
	if !ok {
		t.Fatal("gw1の2つめの後続ジョブがJob型になっていない。")
	}
	if task3.ID() != "task3" {
		t.Error("gw1の1つめの後続ジョブはtask2であるはずが、%sになっている。", task3.ID())
	}
}

func TestSetElements_開始要素と終了要素を設定できる(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 1),
		Flow:  make([]parser.SequenceFlow, 2),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "end"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	if nwk.Start == nil {
		t.Fatal("開始要素が設定されなかった。")
	}
	if nwk.Start.ID() != "task1" {
		t.Error("間違った要素が開始要素に設定された。")
	}

	if nwk.End == nil {
		t.Fatal("終了要素が設定されなかった。")
	}
	if nwk.End.ID() != "task1" {
		t.Error("間違った要素が終了要素に設定された。")
	}
}

func TestSetElements_ジョブ名が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "j1", Name: "jo$b1"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_開始要素の接続先が複数ある場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 2),
		Flow:  make([]parser.SequenceFlow, 2),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "j1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "j2", Name: "job2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "j1"}
	proc.Flow[1] = parser.SequenceFlow{From: "start", To: "j2"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_終了要素の接続元が複数ある場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 2),
		Flow:  make([]parser.SequenceFlow, 2),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "j1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "j2", Name: "job2"}
	proc.Flow[0] = parser.SequenceFlow{From: "j1", To: "end"}
	proc.Flow[1] = parser.SequenceFlow{From: "j2", To: "end"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_フローが空の場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "end"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ID重複時はエラーを吐く_Task同士の重複(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 2),
	}
	proc.Task[0] = parser.ServiceTask{ID: "same", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "same", Name: "job2"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ID重複時はエラーを吐く_Gateway同士の重複(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Gateway: make([]parser.ParallelGateway, 2),
	}
	proc.Gateway[0] = parser.ParallelGateway{ID: "same"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "same"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ID重複時はエラーを吐く_TaskとGatewayの重複(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 1),
		Gateway: make([]parser.ParallelGateway, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "same", Name: "job1"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "same"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_開始要素の接続先が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "noexist"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_終了要素の接続元が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Flow[0] = parser.SequenceFlow{From: "noexist", To: "end"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ジョブの接続先が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 1),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "task1", To: "noexist"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ジョブの接続元が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 1),
		Flow:  make([]parser.SequenceFlow, 1),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "noexist", To: "task1"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ゲートウェイの接続先が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Gateway: make([]parser.ParallelGateway, 1),
		Flow:    make([]parser.SequenceFlow, 1),
	}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Flow[0] = parser.SequenceFlow{From: "gw1", To: "noexist"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ゲートウェイの接続元が不正な場合はエラーを吐く(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Gateway: make([]parser.ParallelGateway, 1),
		Flow:    make([]parser.SequenceFlow, 1),
	}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Flow[0] = parser.SequenceFlow{From: "noexist", To: "gw1"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestSetElements_ジョブの後続が複数ある場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 3),
		Flow:  make([]parser.SequenceFlow, 2),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Flow[0] = parser.SequenceFlow{From: "task1", To: "task2"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "task3"}

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err == nil {
		t.Fatal("エラーが発生しなかった。")
	}
}

func TestLoadJobEx_拡張ジョブファイルが存在しない場合にエラーとしない(t *testing.T) {
	n, _ := NewNetwork("noexistsjobex")
	err := n.LoadJobEx()
	if err != nil {
		t.Errorf("想定外のエラーが発生した: %s", err)
	}
}

func TestSetJobEx_拡張ジョブ情報をセットできる(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 0),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}

	je1 := new(parser.JobEx)
	je1.Node = "node1"
	je1.Port = 1
	je1.FilePath = "path1"
	je1.Param = "param1"
	je1.Env = "env1"
	je1.Workspace = "work1"
	je1.WrnRC = 11
	je1.WrnPtn = "warn1"
	je1.ErrRC = 21
	je1.ErrPtn = "err1"
	je1.TimeoutMin = 60
	je1.SecondaryNode = "secondary"
	je1.SecondaryPort = 2

	je2 := new(parser.JobEx)
	je2.Node = "node2"

	jeMap := make(map[string]*parser.JobEx)
	jeMap["job1"] = je1
	jeMap["job2"] = je2

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	loadTestConfig()
	nwk.setJobEx(jeMap)

	task1 := nwk.elements["task1"].(*Job)
	if task1.Node != "node1" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.Node)
	}
	if task1.Port != 1 {
		t.Errorf("ポート番号[%d]は想定と違っている", task1.Port)
	}
	if task1.FilePath != "path1" {
		t.Errorf("実行ファイル名[%s]は想定と違っている", task1.FilePath)
	}
	if task1.Param != "param1" {
		t.Errorf("実行時パラメータ[%s]は想定と違っている", task1.Param)
	}
	if task1.Env != "env1" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.Env)
	}
	if task1.Workspace != "work1" {
		t.Errorf("作業フォルダ[%s]は想定と違っている", task1.Workspace)
	}
	if task1.WrnRC != 11 {
		t.Errorf("警告条件コード[%d]は想定と違っている", task1.WrnRC)
	}
	if task1.WrnPtn != "warn1" {
		t.Errorf("警告出力パターン[%s]は想定と違っている", task1.WrnPtn)
	}
	if task1.ErrRC != 21 {
		t.Errorf("異常条件コード[%d]は想定と違っている", task1.ErrRC)
	}
	if task1.ErrPtn != "err1" {
		t.Errorf("異常出力パターン[%s]は想定と違っている", task1.ErrPtn)
	}
	if task1.Timeout != 3600 {
		t.Errorf("実行タイムアウト時間[%d]は想定と違っている", task1.Timeout)
	}
	if task1.SecondaryNode != "secondary" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.SecondaryNode)
	}
	if task1.SecondaryPort != 2 {
		t.Errorf("ポート番号[%d]は想定と違っている", task1.SecondaryPort)
	}

	task2 := nwk.elements["task2"].(*Job)
	if task2.Node != "node2" {
		t.Errorf("ノード名[%s]は想定と違っている", task2.Node)
	}
}

func TestSetJobEx_ゼロ値が挿入されたカラムにデフォルト値をセットする(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 0),
	}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}

	je1 := parser.NewJobEx()
	jeMap := make(map[string]*parser.JobEx)
	jeMap["job1"] = je1

	nwk, _ := NewNetwork("test")
	err := nwk.setElements(proc)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s")
	}

	loadTestConfig()
	nwk.setJobEx(jeMap)

	task1 := nwk.elements["task1"].(*Job)
	if task1.Node != "localhost" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.Node)
	}
	if task1.Port != 2015 {
		t.Errorf("ポート番号[%d]は想定と違っている", task1.Port)
	}
	if task1.FilePath != "job1" {
		t.Errorf("実行ファイル名[%s]は想定と違っている", task1.FilePath)
	}
	if task1.Param != "" {
		t.Errorf("実行時パラメータ[%s]は想定と違っている", task1.Param)
	}
	if task1.Env != "" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.Env)
	}
	if task1.Workspace != "" {
		t.Errorf("作業フォルダ[%s]は想定と違っている", task1.Workspace)
	}
	if task1.WrnRC != 0 {
		t.Errorf("警告条件コード[%d]は想定と違っている", task1.WrnRC)
	}
	if task1.WrnPtn != "" {
		t.Errorf("警告出力パターン[%s]は想定と違っている", task1.WrnPtn)
	}
	if task1.ErrRC != 0 {
		t.Errorf("異常条件コード[%d]は想定と違っている", task1.ErrRC)
	}
	if task1.ErrPtn != "" {
		t.Errorf("異常出力パターン[%s]は想定と違っている", task1.ErrPtn)
	}
	if task1.Timeout != 1800 {
		t.Errorf("実行タイムアウト時間[%d]は想定と違っている", task1.Timeout)
	}
	if task1.SecondaryNode != "" {
		t.Errorf("ノード名[%s]は想定と違っている", task1.SecondaryNode)
	}
	if task1.SecondaryPort != 0 {
		t.Errorf("ポート番号[%d]は想定と違っている", task1.SecondaryPort)
	}
}

func TestSetJobEx_ゲートウェイには影響を与えない(t *testing.T) {
	je1 := new(parser.JobEx)
	je1.Node = "node1"
	je1.Port = 1
	je1.FilePath = "path1"
	je1.Param = "param1"
	je1.Env = "env1"
	je1.Workspace = "work1"
	je1.WrnRC = 11
	je1.WrnPtn = "warn1"
	je1.ErrRC = 21
	je1.ErrPtn = "err1"
	je1.TimeoutMin = 60
	je1.SecondaryNode = "secondary"
	je1.SecondaryPort = 2

	jeMap := make(map[string]*parser.JobEx)
	jeMap["job1"] = je1

	n, _ := NewNetwork("test")
	gw1 := NewGateway("job1")
	original := *gw1
	n.elements["job1"] = gw1

	n.setJobEx(jeMap)
	gw1After := n.elements["job1"].(*Gateway)
	if !reflect.DeepEqual(original, *gw1After) {
		t.Error("Gatewayが変更された。")
		t.Log("setJobEx前: %v", original)
		t.Log("setJobEx後: %v", *gw1After)
	}
}

func TestDetectFlowError_正常なフローではエラーを検出しない(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 3),
		Gateway: make([]parser.ParallelGateway, 2),
		Flow:    make([]parser.SequenceFlow, 7),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "gw2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "task3"}
	proc.Flow[4] = parser.SequenceFlow{From: "task2", To: "gw2"}
	proc.Flow[5] = parser.SequenceFlow{From: "task3", To: "gw2"}
	proc.Flow[6] = parser.SequenceFlow{From: "gw2", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err != nil {
		t.Fatalf("想定外のエラーが検出された: %s", err)
	}
}

func TestDetectFlowError_分岐しないゲートウェイは正常なフローとして許容する(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 2),
		Flow:    make([]parser.SequenceFlow, 5),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "gw2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "task2", To: "gw2"}
	proc.Flow[4] = parser.SequenceFlow{From: "gw2", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err != nil {
		t.Fatalf("想定外のエラーが検出された: %s", err)
	}
}

func TestDetectFlowError_startEventがジョブやゲートウェイと接続されていない場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 1),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 1),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "task1", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_endEventがジョブやゲートウェイと接続されていない場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 1),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 1),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_startEventに接続していない始端がある場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 1),
		Flow:    make([]parser.SequenceFlow, 4),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "task2", To: "gw1"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_endEventとの接続要素に別の後続要素がある場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 2),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	task1 := nwk.elements["task1"].(*Job)
	task1.Next = nwk.elements["task2"]

	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_endEventに接続していない終端がある場合はエラー_分岐外(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 2),
		Gateway: make([]parser.ParallelGateway, 0),
		Flow:    make([]parser.SequenceFlow, 2),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task2", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_endEventに接続していない終端がある場合はエラー_分岐内(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 3),
		Gateway: make([]parser.ParallelGateway, 2),
		Flow:    make([]parser.SequenceFlow, 6),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "gw2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "task3"}
	proc.Flow[4] = parser.SequenceFlow{From: "task2", To: "gw2"}
	proc.Flow[5] = parser.SequenceFlow{From: "gw2", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_分岐を結合せずに終了した場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 3),
		Gateway: make([]parser.ParallelGateway, 1),
		Flow:    make([]parser.SequenceFlow, 6),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "task3"}
	proc.Flow[4] = parser.SequenceFlow{From: "task2", To: "end"}
	proc.Flow[5] = parser.SequenceFlow{From: "task3", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_分岐がネストした場合はエラー(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 4),
		Gateway: make([]parser.ParallelGateway, 3),
		Flow:    make([]parser.SequenceFlow, 10),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Task[3] = parser.ServiceTask{ID: "task4", Name: "job4"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "gw2"}
	proc.Gateway[2] = parser.ParallelGateway{ID: "gw3"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "gw2"}
	proc.Flow[4] = parser.SequenceFlow{From: "gw2", To: "task3"}
	proc.Flow[5] = parser.SequenceFlow{From: "gw2", To: "task4"}
	proc.Flow[6] = parser.SequenceFlow{From: "task2", To: "gw3"}
	proc.Flow[7] = parser.SequenceFlow{From: "task3", To: "gw3"}
	proc.Flow[8] = parser.SequenceFlow{From: "task4", To: "gw3"}
	proc.Flow[9] = parser.SequenceFlow{From: "gw3", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)
	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_定義外の型の要素が発見されたらエラー_分岐外(t *testing.T) {
	proc := &parser.Process{
		Start: make([]parser.StartEvent, 1),
		End:   make([]parser.EndEvent, 1),
		Task:  make([]parser.ServiceTask, 1),
		Flow:  make([]parser.SequenceFlow, 3),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)

	task1 := nwk.elements["task1"].(*Job)
	undef := new(undefinedElement)
	nwk.elements[undef.ID()] = undef
	undef.AddNext(task1)
	nwk.Start = undef

	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestDetectFlowError_定義外の型の要素が発見されたらエラー_分岐内(t *testing.T) {
	proc := &parser.Process{
		Start:   make([]parser.StartEvent, 1),
		End:     make([]parser.EndEvent, 1),
		Task:    make([]parser.ServiceTask, 3),
		Gateway: make([]parser.ParallelGateway, 2),
		Flow:    make([]parser.SequenceFlow, 7),
	}
	proc.Start[0] = parser.StartEvent{ID: "start"}
	proc.End[0] = parser.EndEvent{ID: "end"}
	proc.Task[0] = parser.ServiceTask{ID: "task1", Name: "job1"}
	proc.Task[1] = parser.ServiceTask{ID: "task2", Name: "job2"}
	proc.Task[2] = parser.ServiceTask{ID: "task3", Name: "job3"}
	proc.Gateway[0] = parser.ParallelGateway{ID: "gw1"}
	proc.Gateway[1] = parser.ParallelGateway{ID: "gw2"}
	proc.Flow[0] = parser.SequenceFlow{From: "start", To: "task1"}
	proc.Flow[1] = parser.SequenceFlow{From: "task1", To: "gw1"}
	proc.Flow[2] = parser.SequenceFlow{From: "gw1", To: "task2"}
	proc.Flow[3] = parser.SequenceFlow{From: "gw1", To: "task3"}
	proc.Flow[4] = parser.SequenceFlow{From: "task2", To: "gw2"}
	proc.Flow[5] = parser.SequenceFlow{From: "task3", To: "gw2"}
	proc.Flow[6] = parser.SequenceFlow{From: "gw2", To: "end"}

	nwk, _ := NewNetwork("test")
	nwk.setElements(proc)

	task2 := nwk.elements["task2"].(*Job)
	undef := new(undefinedElement)

	nwk.elements[undef.ID()] = undef
	undef.AddNext(task2.Next)
	task2.Next = undef

	err := nwk.DetectFlowError()
	if err == nil {
		t.Fatalf("エラーが検出されていない。")
	}
}

func TestNetworkRun_ネットワークを実行できる(t *testing.T) {
	loadTestConfig()

	n, _ := NewNetwork("test")
	j1 := generateTestJob(1)
	j2 := generateTestJob(2)
	j3 := generateTestJob(3)
	j4 := generateTestJob(4)
	g1 := NewGateway("gwid1")
	g2 := NewGateway("gwid2")
	n.elements[j1.ID()] = j1
	n.elements[j2.ID()] = j2
	n.elements[j3.ID()] = j3
	n.elements[j4.ID()] = j4
	n.elements[g1.ID()] = g1
	n.elements[g2.ID()] = g2

	n.Start = j1
	j1.AddNext(g1)
	g1.AddNext(j2)
	g1.AddNext(j3)
	j2.AddNext(g2)
	j3.AddNext(g2)
	g2.AddNext(j4)
	n.End = j4

	if err := n.Run(); err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if !j1.isExecuted {
		t.Errorf("j1が実行されていない。")
	}
	if !j2.isExecuted {
		t.Errorf("j2が実行されていない。")
	}
	if !j3.isExecuted {
		t.Errorf("j3が実行されていない。")
	}
	if !j4.isExecuted {
		t.Errorf("j4が実行されていない。")
	}
}

func TestNetworkRun_開始要素がnilの場合はエラー(t *testing.T) {
	loadTestConfig()

	n, _ := NewNetwork("test")
	j1 := generateTestJob(1)
	j1.hasError = true
	n.elements[j1.ID()] = j1

	n.End = j1

	if err := n.Run(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestNetworkRun_ジョブネットワークの開始処理で失敗したらエラー(t *testing.T) {
	config.DB.DBFile = ""

	n, _ := NewNetwork("test")
	j1 := generateTestJob(1)
	n.elements[j1.ID()] = j1

	n.Start = j1
	n.End = j1

	if err := n.Run(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
	if j1.isExecuted {
		t.Error("j1が実行された。")
	}
}

func TestNetworkRun_ネットワーク中のジョブが異常終了したらエラー(t *testing.T) {
	loadTestConfig()

	n, _ := NewNetwork("test")
	j1 := generateTestJob(1)
	j1.hasError = true
	n.elements[j1.ID()] = j1

	n.Start = j1
	n.End = j1

	if err := n.Run(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestNetworkRun_ジョブネットワークがendEvent以外で終端したらエラー(t *testing.T) {
	loadTestConfig()

	n, _ := NewNetwork("test")
	j1 := generateTestJob(1)
	n.elements[j1.ID()] = j1

	n.Start = j1

	if err := n.Run(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}
