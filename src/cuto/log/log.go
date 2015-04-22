// Copyright 2015 unirita Inc.
// Created 2015/04/15 honda

package log

import (
	"fmt"
	"os"

	"github.com/cihub/seelog"

	"cuto/util"
)

const lockTimeout = 1000
const mutexHeader = "Unirita_CuteLog_"

var valid = false
var locker *util.MutexHandle

// ロガーの初期化処理を行う
//
// param ; dir       ログファイルの出力先ディレクトリ。
//
// param : name      ログファイルの種別（例：master、servant）。
//
// param : level     出力ログレベル（trace,debug,info,warn,error,criticalのいずれかを指定）
//
// param : maxSizeKB ログファイルの最大サイズ。この値を超えるとログローテーションが発生する。
//
// param : maxRolls  ログファイルの最大世代数
//
// return : エラー情報を返す。
func Init(dir string, name string, level string, maxSizeKB int, maxRolls int) error {
	config := generateConfigString(dir, name, level, maxSizeKB, maxRolls)
	logger, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		return err
	}

	mutexName := mutexHeader + name
	locker, err = util.InitMutex(mutexName)
	if err != nil {
		return err
	}

	seelog.ReplaceLogger(logger)
	valid = true

	return nil
}

// ロガーの終了処理を行う。
func Term() {
	locker.TermMutex()
	valid = false
}

// traceレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Trace(msg ...interface{}) {
	if !valid {
		return
	}
	seelog.Trace(msg...)
	Flush()
}

// debugレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Debug(msg ...interface{}) {
	if !valid {
		return
	}
	seelog.Debug(msg...)
	Flush()
}

// infoレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Info(msg ...interface{}) {
	if !valid {
		return
	}
	seelog.Info(msg...)
	Flush()
}

// warnレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Warn(msg ...interface{}) {
	if !valid {
		return
	}
	seelog.Warn(msg...)
	Flush()
}

// errorレベルでログメッセージを出力する。
//
// param : msg 出力するメッセージ。複数指定した場合は結合して出力される。
func Error(msg ...interface{}) {
	if !valid {
		return
	}
	seelog.Error(msg...)
	Flush()
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

// ログのフラッシュを行う。
func Flush() {
	if !valid {
		return
	}
	locker.Lock(lockTimeout)
	defer locker.Unlock()
	seelog.Flush()
}

func generateConfigString(dir string, name string, level string, maxSizeKB int, maxRolls int) string {
	format := `
<seelog minlevel="%s">
    <outputs formatid="common">
        <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
    </outputs>
    <formats>
        <format id="common" format="%%Date(2006-01-02 15:04:05.000) [%%LEV] %%Msg%%n"/>
    </formats>
</seelog>`

	filename := fmt.Sprintf("%s%c%s.log", dir, os.PathSeparator, name)

	// rollingfileのmaxrollsの数字は、書き込み中のログファイルを含まずにカウントするため引数をデクリメントする。
	maxRolls--
	return fmt.Sprintf(format, level, filename, maxSizeKB*1024, maxRolls)
}
