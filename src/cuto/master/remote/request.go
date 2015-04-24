// Copyright 2015 unirita Inc.
// Created 2015/04/10 honda

package remote

import (
	"fmt"
	"net"
	"strings"
	"time"

	"cuto/log"
	"cuto/master/config"
	"cuto/message"
)

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
	_, err = conn.Write([]byte(req))
	if err != nil {
		return ``, err
	}

	buf := make([]byte, bufSize)
	var res string

WAITRESPONSE:
	for {
		select {
		case err = <-readResponse(conn, &buf):
			if err != nil {
				return ``, err
			}
			res = string(buf)
			log.Debug(res)
			if res == message.HEARTBEAT {
				// ハートビートメッセージの場合はバッファサイズを初期化して再度read待ちをする。
				//@todo 1度目の受信でジョブの起動ステータスを更新したい。
				buf = buf[:bufSize]
				continue
			} else if strings.HasPrefix(res, message.ST_HEADER) {
				st := res[len(message.ST_HEADER):]
				stCh <- st
				buf = buf[:bufSize]
				continue
			}

			break WAITRESPONSE
		case <-time.After(timeout):
			return ``, fmt.Errorf("Connetion timeout.")
		}
	}

	return res, nil
}

func readResponse(c net.Conn, b *[]byte) <-chan error {
	ch := make(chan error, 10)
	go func() {
		l, err := c.Read(*b)
		if err != nil {
			ch <- err
			return
		}

		*b = (*b)[:l]
		ch <- nil
	}()

	return ch
}
