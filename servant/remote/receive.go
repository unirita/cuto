// マスタからの受信
// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package remote

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/unirita/cuto/console"
	"github.com/unirita/cuto/log"
)

// 送受信メッセージの終端文字
const MsgEnd = "\n"

const protocol = "tcp"

// 引数で指定したbind用のアドレスとポート番号portを指定してメッセージの受信待ちを開始する。
// 受信したメッセージは戻り値のchan stringに順番に挿入されていく。
// 引数で指定した多重度でチャネルを作成する。
//
// 引数：bindAddr バインドアドレス
//
// 引数：port Listenポート番号
//
// 引数：multi メッセージ受信の多重度
//
// 戻り値：メッセージ受信を通知する、受信チャネル
//
// 戻り値：エラー情報
func StartReceive(bindAddr string, port int, multi int) (<-chan *Session, error) {
	addr := fmt.Sprintf("%s:%d", bindAddr, port)

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}
	console.Display("CTS007I", bindAddr, port)

	sq := make(chan *Session, multi)
	go receiveLoop(listener, sq)

	return sq, nil
}

func receiveLoop(listener net.Listener, sq chan<- *Session) {
	for {
		receiveLoopProcess(listener, sq)
	}
}

func receiveLoopProcess(listener net.Listener, sq chan<- *Session) error {
	conn, err := listener.Accept()
	if err != nil {
		log.Error(err)
		return err
	}

	console.Display("CTS014I")
	go receiveMessage(conn, sq)
	return nil
}

func receiveMessage(conn net.Conn, sq chan<- *Session) error {
	const timeout = 10

	deadLine := time.Now().Add(timeout * time.Second)
	if err := conn.SetReadDeadline(deadLine); err != nil {
		log.Error(fmt.Sprintf("[%v]: %v\n", conn.RemoteAddr(), err))
		return err
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		log.Error(fmt.Sprintf("[%v]: %v\n", conn.RemoteAddr(), scanner.Err()))
		return scanner.Err()
	}
	sq <- NewSession(conn, scanner.Text())
	return nil
}
