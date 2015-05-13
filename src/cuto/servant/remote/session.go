// マスタからの接続要求を受け取る。
// 作成者：2015/04/09　本田
// copyright. unirita Inc.

package remote

import (
	"net"
	"time"

	"cuto/console"
	"cuto/db"
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

	endHeartbeatCh chan endSig
	doJobRequest   func(req *message.Request, conf *config.ServantConfig, stCh chan<- string) *message.Response
}

// Sessionオブジェクトのコンストラクタ
func NewSession(conn net.Conn, body string) *Session {
	s := new(Session)
	s.Conn = conn
	s.Body = body
	s.doJobRequest = job.DoJobRequest
	s.startHeartbeat()
	return s
}

// セッションに対応したジョブ実行要求に基いてジョブを実行する。
// 引数：conf 設定情報
// 戻り値：なし
func (s *Session) Do(conf *config.ServantConfig) error {
	defer s.Conn.Close()
	defer s.endHeartbeat()

	req := new(message.Request)
	if err := req.ParseJSON(s.Body); err != nil {
		console.Display("CTS015E", err.Error())
		return err
	}

	res, err := s.doRequest(req, conf)
	if err != nil {
		log.Error(err)
		res = s.createErrorResponse(req, err)
	}

	var resMsg string
	resMsg, err = res.GenerateJSON()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = s.Conn.Write([]byte(resMsg + MsgEnd))
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *Session) doRequest(req *message.Request, conf *config.ServantConfig) (*message.Response, error) {
	err := req.ExpandServantVars()
	if err != nil {
		console.Display("CTS015E", err.Error())
		return nil, err
	}

	stCh := make(chan string, 1)
	go s.waitAndSendStartTime(stCh)
	defer close(stCh)

	res := s.doJobRequest(req, conf, stCh)
	return res, nil
}

// ハートビートを開始する。
func (s *Session) startHeartbeat() {
	s.endHeartbeatCh = make(chan endSig, 1)
	go func() {
		t := time.Duration(config.Servant.Job.HeartbeatSpanSec) * time.Second
		for {
			select {
			case <-s.endHeartbeatCh:
				return
			case <-time.After(t):
				log.Debug("send heatbeat...")
				s.Conn.Write([]byte(message.HEARTBEAT + MsgEnd))
			}
		}
	}()
}

// ハートビートメッセージを停止する
func (s *Session) endHeartbeat() {
	s.endHeartbeatCh <- endSig{}
	close(s.endHeartbeatCh)
}

// スタート時刻の決定を待ち、masterへ送信する。
func (s *Session) waitAndSendStartTime(stCh <-chan string) {
	st := <-stCh
	if len(st) == 0 {
		// （主にチャネルがクローズされることにより）空文字列が送られてきた場合は何もしない。
		return
	}

	s.Conn.Write([]byte(message.ST_HEADER + st + MsgEnd))
}

// ジョブ実行結果が得られないようなエラーが発生した場合のレスポンスメッセージを生成する。
func (s *Session) createErrorResponse(req *message.Request, err error) *message.Response {
	res := new(message.Response)
	res.NID = req.NID
	res.JID = req.JID
	res.Stat = db.ABNORMAL
	res.Detail = err.Error()
	return res
}
