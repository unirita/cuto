// マスタからの受信
// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package remote

import (
	"fmt"
	"net"
	"time"

	"cuto/console"
	"cuto/log"
)

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

	log.Debug("addr = ", addr)

	listener, err := net.Listen(protocol, addr)
	if err != nil {
		return nil, err
	}

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
	const bufSize = 1024

	var readLen int
	var err error

	err = conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
	if err != nil {
		log.Error(fmt.Sprintf("[%v]: %v\n", conn.RemoteAddr(), err))
		return err
	}

	buf := make([]byte, bufSize)
	readLen, err = conn.Read(buf)
	if err != nil {
		log.Error(fmt.Sprintf("[%v]: %v\n", conn.RemoteAddr(), err))
		return err
	}

	sq <- NewSession(conn, string(buf[:readLen]))
	return nil
}
