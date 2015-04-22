// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package main

import (
	"syscall"

	"os"
	"os/signal"

	"cuto/log"
	"cuto/servant/config"
	"cuto/servant/remote"
)

// サーバントメインルーチン
func Run() (int, error) {
	// セッションの用意
	sq, err := remote.StartReceive(config.Servant.Sys.BindAddress, config.Servant.Sys.BindPort, config.Servant.Job.MultiProc)
	if err != nil {
		return -1, err
	}

	// OSシグナル受信用チャネル、SIGTERMとSIGHUPに対応する
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

LOOP:
	for {
		select {
		case ch := <-signalCh: // OSからシグナルキャッチ
			if isEnd := isEndSig(ch); isEnd {
				break LOOP
			}
		case session := <-sq: // マスタからの要求受信
			go session.Do(config.Servant)
		}
	}
	return 0, nil
}

func isEndSig(sig os.Signal) bool {
	log.Debug("Receive Signal : ", sig)
	if sig == syscall.SIGTERM || sig == syscall.SIGHUP {
		// ハングアップ？
		config.ReloadConfig()
	} else if sig == syscall.SIGINT {
		return true
	}
	return false
}
