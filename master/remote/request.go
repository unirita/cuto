// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package remote

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/unirita/cuto/log"
	"github.com/unirita/cuto/master/config"
	"github.com/unirita/cuto/message"
)

type response struct {
	msg string
	err error
}

// 送受信メッセージの終端文字
const MsgEnd = "\n"

// ホスト名がhost、ポート番号がportのservantへ接続し、ジョブ実行要求を送信する。
// servantから返信されたジョブ実行結果を関数外へ返す。
//
// param : host ホスト名。
//
// param : port ポート番号。
//
// param : req リクエストメッセージ。
//
// return : 返信メッセージ。
//
// return : エラー情報。
func SendRequest(host string, port int, req string, stCh chan<- string) (string, error) {
	const bufSize = 1024
	timeout := time.Duration(config.Job.ConnectionTimeoutSec) * time.Second

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return ``, err
	}
	defer conn.Close()

	log.Debug(req)
	_, err = conn.Write([]byte(req + MsgEnd))
	if err != nil {
		return ``, err
	}

	scanner := bufio.NewScanner(conn)

	for {
		select {
		case res := <-readResponse(scanner):
			if res.err != nil {
				return ``, res.err
			}
			log.Debug(res.msg)
			if res.msg == message.HEARTBEAT {
				// ハートビートメッセージの場合はバッファサイズを初期化する。
				continue
			} else if strings.HasPrefix(res.msg, message.ST_HEADER) {
				st := res.msg[len(message.ST_HEADER):]
				stCh <- st
				continue
			}

			return res.msg, nil
		case <-time.After(timeout):
			return ``, fmt.Errorf("Connetion timeout.")
		}
	}
}

func readResponse(scanner *bufio.Scanner) <-chan *response {
	ch := make(chan *response, 10)
	go func() {
		res := new(response)
		if scanner.Scan() {
			res.msg = scanner.Text()
		} else {
			res.err = scanner.Err()
		}
		ch <- res
		if res.err != nil {
			return
		}
	}()

	return ch
}
