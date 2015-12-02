// Copyright 2015 unirita Inc.
// Created 2015/04/15 honda

package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cihub/seelog"

	"github.com/unirita/cuto/util"
)

const lockHeader = "Unirita_CutoLog_"

var lockTimeout = 1000
var valid = false
var locker *util.LockHandle

// ロガーの初期化処理を行う
//
// param ; dir       ログファイルの出力先ディレクトリ。
//
// param : name      ログファイルの種別（例：master、servant）。
//
// param : identifer ロック用ファイルに（付与する識別ID（例：servantはListePort）。
//
// param : level     出力ログレベル（trace,debug,info,warn,error,criticalのいずれかを指定）
//
// param : maxSizeKB ログファイルの最大サイズ。この値を超えるとログローテーションが発生する。
//
// param : maxRolls  ログファイルの最大世代数
//
// param : timeoutSec  ロックのタイムアウト秒
//
// return : エラー情報を返す。
func Init(dir string, name string, identifer string, level string, maxSizeKB int, maxRolls int, timeoutSec int) error {
	var lockErr error
	lockName := lockHeader + name
	if identifer != "" {
		lockName = lockName + "_" + identifer
	}
	lockName = lockName + ".lock"
	locker, lockErr = util.InitLock(lockName)
	if lockErr != nil {
		return lockErr
	}
	if timeoutSec > 0 {
		lockTimeout = timeoutSec * 1000
	}

	logfile := filepath.Join(dir, name) + ".log"
	config := generateConfigString(logfile, level, maxSizeKB, maxRolls)
	logger, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		Term()
		return err
	}

	if err := makeFileIfNotExist(logfile); err != nil {
		Term()
		return err
	}

	seelog.ReplaceLogger(logger)
	valid = true

	return nil
}

// ロガーの終了処理を行う。
func Term() {
	locker.TermLock()
	valid = false
}

// traceレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Trace(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Trace(msg...)
}

// debugレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Debug(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Debug(msg...)
}

// infoレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Info(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Info(msg...)
}

// warnレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Warn(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Warn(msg...)
}

// errorレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Error(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Error(msg...)
}

// criticalレベルでログメッセージを出力する。
// この関数が呼び出されると、ただちにログのフラッシュが行われる。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Critical(msg ...interface{}) {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Critical(msg...)
}

func makeFileIfNotExist(logfile string) error {
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	if _, err := os.Stat(logfile); !os.IsNotExist(err) {
		// ファイルが存在する場合は何もしない。
		// os.IsExistはerr=nilのときfalseを返すため、os.IsNotExistで判定している。
		return nil
	}

	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}

func generateConfigString(logfile string, level string, maxSizeKB int, maxRolls int) string {
	format := `
<seelog type="sync" minlevel="%s">
    <outputs formatid="common">
        <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
    </outputs>
    <formats>
        <format id="common" format="%%Date(2006-01-02 15:04:05.000) [%d] [%%LEV] %%Msg%%n"/>
    </formats>
</seelog>`

	// rollingfileのmaxrollsの数字は、書き込み中のログファイルを含まずにカウントするため引数をデクリメントする。
	maxRolls--
	return fmt.Sprintf(format, level, logfile, maxSizeKB*1024, maxRolls, os.Getpid())
}
