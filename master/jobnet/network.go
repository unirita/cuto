// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package jobnet

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/unirita/cuto/console"
	"github.com/unirita/cuto/db"
	"github.com/unirita/cuto/db/tx"
	"github.com/unirita/cuto/log"
	"github.com/unirita/cuto/master/config"
	"github.com/unirita/cuto/master/jobnet/parser"
	"github.com/unirita/cuto/message"
	"github.com/unirita/cuto/util"
)

// ジョブネット全体を表す構造体
type Network struct {
	ID         int                // ジョブネットワークID。
	Name       string             // ジョブネットワーク名。
	Start      Element            // スタートイベントのノード。
	End        Element            // エンドイベントのノード。
	MasterPath string             // ジョブネットワークファイルパス。
	JobExPath  string             // 拡張ジョブ定義ファイルパス。
	elements   map[string]Element // ジョブネットワークの構成要素Map。
	Result     *tx.ResultMap      // 実行結果情報。
	globalLock *util.LockHandle   // マスタ間ロックハンドル
	localMutex sync.Mutex         // ゴルーチン間のミューテックス
}

// cuto masterが使用するミューテックス名。
const lock_name string = "Unirita_CutoMaster.lock"

// Network構造体のコンストラクタ関数
//
// param : name ジョブネットワーク名。
//
// return : ジョブネットワーク構造体。
func NewNetwork(name string) (*Network, error) {
	nwk := new(Network)
	nwk.Name = name
	nwk.elements = make(map[string]Element)
	filePrefix := filepath.Join(config.Dir.JobnetDir, name)
	nwk.MasterPath = filePrefix + ".bpmn"
	nwk.JobExPath = filePrefix + ".csv"

	var err error
	nwk.globalLock, err = util.InitLock(lock_name)
	if err != nil {
		return nil, err
	}
	return nwk, err
}

// ネットワーク名nameを元にネットワーク定義ファイルをロードし、Network構造体のオブジェクトを返す。
//
// param : name ジョブネットワーク名。
//
// return : ジョブネットワーク構造体。
func LoadNetwork(name string) *Network {
	nwk, err := NewNetwork(name)
	if err != nil {
		console.Display("CTM019E", err)
		return nil
	}

	file, err := os.Open(nwk.MasterPath)
	if err != nil {
		console.Display("CTM010E", nwk.MasterPath)
		log.Error(err)
		return nil
	}
	defer file.Close()

	err = nwk.LoadElements(file)
	if err != nil {
		console.Display("CTM011E", nwk.MasterPath, err)
		return nil
	}

	return nwk
}

// io.Readerからネットワーク定義を読み込み、n.Start/n.End/n.elementsへ値をセットする。
//
// param : r Reader。
//
// return : エラー情報。
func (n *Network) LoadElements(r io.Reader) error {
	proc, err := parser.ParseNetwork(r)
	if err != nil {
		return err
	}

	return n.setElements(proc)
}

// BPMNパース結果のProcess構造体からネットワークの各要素を取得し、セットする。
func (n *Network) setElements(proc *parser.Process) error {
	for _, t := range proc.Task {
		if _, exists := n.elements[t.ID]; exists {
			return fmt.Errorf("Element[id = %s] duplicated.", t.ID)
		}

		var err error
		n.elements[t.ID], err = NewJob(t.ID, t.Name, n)
		if err != nil {
			return err
		}
	}

	for _, g := range proc.Gateway {
		if _, exists := n.elements[g.ID]; exists {
			return fmt.Errorf("Element[id = %s] duplicated.", g.ID)
		}
		n.elements[g.ID] = NewGateway(g.ID)
	}

	sid := proc.Start[0].ID
	eid := proc.End[0].ID

	for _, f := range proc.Flow {
		switch {
		case f.From == sid:
			if n.Start != nil {
				return fmt.Errorf("StartEvent cannot connect with over 1 element.")
			}

			if f.To == eid {
				return fmt.Errorf("Jobnet is empty.")
			}

			var ok bool
			n.Start, ok = n.elements[f.To]
			if !ok {
				return fmt.Errorf("StartEvent connects with imaginary element[id = %s].", f.To)
			}
		case f.To == eid:
			if n.End != nil {
				return fmt.Errorf("EndEvent cannot connect with over 1 element.")
			}

			var ok bool
			n.End, ok = n.elements[f.From]
			if !ok {
				return fmt.Errorf("EndEvent connects with imaginary element[id = %s].", f.From)
			}
		default:
			from, ok := n.elements[f.From]
			if !ok {
				return fmt.Errorf("There is a sequenceFlow which refers imaginary element[id = %s].", f.From)
			}
			to, ok := n.elements[f.To]
			if !ok {
				return fmt.Errorf("There is a sequenceFlow which refers imaginary element[id = %s].", f.To)
			}
			if err := from.AddNext(to); err != nil {
				return err
			}
		}
	}

	return nil
}

// JobExファイルをロードし、ネットワーク内のジョブへ拡張ジョブ定義をセットする。
//
// return : エラー情報。
func (n *Network) LoadJobEx() error {
	jobEx, err := parser.ParseJobExFile(n.JobExPath)
	if err != nil {
		return err
	}
	n.setJobEx(jobEx)

	return nil
}

// ネットワーク内のジョブへ拡張ジョブ定義のパース結果をセットする。
func (n *Network) setJobEx(m map[string]*parser.JobEx) {
	for _, e := range n.elements {
		switch e.(type) {
		case *Job:
			j := e.(*Job)
			if je, ok := m[j.Name]; ok {
				j.Node = je.Node
				j.Port = je.Port
				j.FilePath = je.FilePath
				j.Param = je.Param
				j.Env = je.Env
				j.Workspace = je.Workspace
				j.WrnRC = je.WrnRC
				j.WrnPtn = je.WrnPtn
				j.ErrRC = je.ErrRC
				j.ErrPtn = je.ErrPtn
				j.Timeout = je.TimeoutMin * 60
				j.SecondaryNode = je.SecondaryNode
				j.SecondaryPort = je.SecondaryPort
			}
			j.SetDefaultEx()
		default:
			continue
		}
	}
}

// 実行フローのエラー検出を行う。
//
// return : エラー情報。
func (n *Network) DetectFlowError() error {
	if n.Start == nil {
		return fmt.Errorf("There is no element which connects with startEvent.")
	}
	if n.End == nil {
		return fmt.Errorf("There is no element which connects with endEvent.")
	}

	novisit := make(map[string]Element)
	for k, v := range n.elements {
		novisit[k] = v
	}

	err := n.scanFlow(n.Start, novisit)
	if err != nil {
		return err
	}

	if len(novisit) > 0 {
		return fmt.Errorf("Isolated element is detected.")
	}

	return nil
}

func (n *Network) scanFlow(e Element, novisit map[string]Element) error {
	delete(novisit, e.ID())
	if e == n.End {
		if e.HasNext() {
			return fmt.Errorf("Element which connects with endEvent cannot connect with another element.")
		}
		return nil
	} else if !e.HasNext() {
		return fmt.Errorf("Element[id = %s] cannot terminate network because it is not a endEvent.", e.ID)
	}

	switch e.(type) {
	case *Job:
		j := e.(*Job)
		return n.scanFlow(j.Next, novisit)
	case *Gateway:
		g := e.(*Gateway)

		if len(g.Nexts) == 1 {
			return n.scanFlow(g.Nexts[0], novisit)
		} else {
			var jct Element = nil
			for _, branch := range g.Nexts {
				bind, err := n.scanFlowParallel(branch, novisit)
				if err != nil {
					return err
				}

				if jct == nil {
					jct = bind
				} else if jct != bind {
					return fmt.Errorf("Cannot nest branches.")
				}
			}

			return n.scanFlow(jct, novisit)
		}
	default:
		return fmt.Errorf("Irregal element was detected.")
	}
}

func (n *Network) scanFlowParallel(e Element, novisit map[string]Element) (Element, error) {
	delete(novisit, e.ID())
	switch e.(type) {
	case *Job:
		j := e.(*Job)
		if j.Next == nil {
			return nil, fmt.Errorf("EndEvent cannot connect with branch.")
		}
		return n.scanFlowParallel(j.Next, novisit)
	case *Gateway:
		return e, nil
	default:
		return nil, fmt.Errorf("Irregal element was detected.")
	}
}

// ネットワークを実行する。
//
// return : エラー情報。
func (n *Network) Run() error {
	if n.Start == nil {
		return fmt.Errorf("Start element of network is nil.")
	}

	err := n.start()
	if err != nil {
		console.Display("CTM019E", err)
		return err
	}
	console.Display("CTM012I", n.Name, n.ID)

	return n.runNodes()
}

// ネットワークをリランする。
//
// return : エラー情報。
func (n *Network) Rerun() error {
	if n.Start == nil {
		return fmt.Errorf("Start element of network is nil.")
	}

	err := n.resume()
	if err != nil {
		console.Display("CTM019E", err)
		return err
	}

	prePID := n.Result.JobnetResult.PID
	if util.IsProcessExists(prePID) {
		return fmt.Errorf("JOBNETWORK [%d] still running.", n.ID)
	}

	console.Display("CTM012I", n.Name, n.ID)
	n.setIsRerunJob()

	return n.runNodes()
}

func (n *Network) runNodes() error {
	current := n.Start
	for {
		next, err := current.Execute()
		if err != nil {
			return n.end(err)
		}
		if current == n.End {
			return n.end(nil)
		} else if next == nil {
			err := fmt.Errorf("Element[id = %s] cannot terminate network because it is not a endEvent.", current.ID())
			return n.end(err)
		}
		current = next
	}
	panic("Not reached.")
}

func (n *Network) setIsRerunJob() {
	for _, e := range n.elements {
		if j, ok := e.(*Job); ok {
			if _, exists := n.Result.Jobresults[j.ID()]; exists {
				j.IsRerunJob = true
			}
		}
	}
}

// ジョブネットワークの開始処理
func (n *Network) start() error {
	timeout := config.Job.DefaultTimeoutMin * 60 * 1000
	if timeout <= 0 {
		timeout = 60000
	}

	err := n.globalLock.Lock(timeout)
	if err != nil {
		if err != util.ErrBusy {
			return err
		}
		return fmt.Errorf("Lock Timeout.")
	}
	defer n.globalLock.Unlock()

	n.Result, err = tx.StartJobNetwork(n.Name, config.DB.DBFile)
	if err != nil {
		return err
	}

	n.ID = n.Result.JobnetResult.ID
	message.AddSysValue(`JOBNET`, `ID`, strconv.Itoa(n.ID))
	message.AddSysValue(`JOBNET`, `SD`, n.Result.JobnetResult.StartDate)

	return nil
}

// ジョブネットワークの再実行開始処理
func (n *Network) resume() error {
	timeout := config.Job.DefaultTimeoutMin * 60 * 1000
	if timeout <= 0 {
		timeout = 60000
	}

	err := n.globalLock.Lock(timeout)
	if err != nil {
		if err != util.ErrBusy {
			return err
		}
		return fmt.Errorf("Lock Timeout.")
	}
	defer n.globalLock.Unlock()

	n.Result, err = tx.ResumeJobNetwork(n.ID, config.DB.DBFile)
	if err != nil {
		return err
	}

	message.AddSysValue(`JOBNET`, `ID`, strconv.Itoa(n.ID))
	message.AddSysValue(`JOBNET`, `SD`, n.Result.JobnetResult.StartDate)

	return nil
}

// ジョブネットワークの終了処理
func (n *Network) end(err error) error {
	if err != nil {
		n.Result.EndJobNetwork(db.ABNORMAL, err.Error())
	} else {
		n.Result.EndJobNetwork(db.NORMAL, "")
	}
	return err
}

// 終了処理を行う。
func (n *Network) Terminate() {
	n.globalLock.TermLock()
}
