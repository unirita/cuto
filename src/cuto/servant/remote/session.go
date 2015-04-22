// マスタからの接続要求を受け取る。
// 作成者：2015/04/09　本田
// copyright. unirita Inc.

package remote

import (
	"net"
	"time"

	"cuto/console"
	"cuto/log"
	"cuto/message"
	"cuto/servant/config"
	"cuto/servant/job"
)

type endSig struct{}

// masterからのジョブ実行要求1つに対応する構造体。
type Session struct {
	Conn net.Conn
	Body string

	doRequest func(req *message.Request, conf *config.ServantConfig, stCh chan<- string) *message.Response
}

// Sessionオブジェクトのコンストラクタ
func NewSession(conn net.Conn, body string) *Session {
	return &Session{Conn: conn, Body: body, doRequest: job.DoJobRequest}
}

// セッションに対応したジョブ実行要求に基いてジョブを実行する。
// 引数：conf 設定情報
// 戻り値：なし
func (s *Session) Do(conf *config.ServantConfig) {
	defer s.Conn.Close()

	var err error
	req := new(message.Request)
	err = req.ParseJSON(s.Body)
	if err != nil {
		console.Display("CTS015E", err.Error())
	}

	err = req.ExpandServantVars()
	if err != nil {
		console.Display("CTS015E", err.Error())
	}

	stCh := make(chan string, 1)
	go s.waitAndSendStartTime(stCh)

	endCh := s.startHeartBeat()
	res := s.doRequest(req, conf, stCh)
	endCh <- endSig{}

	close(stCh)

	var resMsg string
	resMsg, err = res.GenerateJSON()
	if err != nil {
		log.Error(err)
	}

	_, err = s.Conn.Write([]byte(resMsg))
	if err != nil {
		log.Error(err)
	}
}

// ハートビートを開始する。
func (s *Session) startHeartBeat() chan endSig {
	ch := make(chan endSig, 1)

	go func() {
		t := time.Duration(config.Servant.Job.HeartbeatSpanSec) * time.Second

	HEARTBEATLOOP:
		for {
			select {
			case <-ch:
				break HEARTBEATLOOP
			case <-time.After(t):
				log.Debug("send heatbeat...")
				s.Conn.Write([]byte(message.HEARTBEAT))
			}
		}
	}()

	return ch
}

// スタート時刻の決定を待ち、masterへ送信する。
func (s *Session) waitAndSendStartTime(stCh <-chan string) {
	st := <-stCh
	if len(st) == 0 {
		// （主にチャネルがクローズされることにより）空文字列が送られてきた場合は何もしない。
		return
	}

	s.Conn.Write([]byte(message.ST_HEADER + st))
}
