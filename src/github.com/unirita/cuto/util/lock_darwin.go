// Copyright 2015 unirita Inc.
// Created 2015/06/03 shanxia

package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type LockHandle struct {
	fd     int  //lockファイルディスクリプタ
	isLock bool // ロックフラグ
}

var (
	ErrBusy = errors.New("Locked by other process.")

	lockFilePath = os.TempDir()
)

// ファイルを利用した同期処理機能の初期化関数。
// ファイル作成が可能なファイル名を指定します。
func InitLock(name string) (*LockHandle, error) {
	if len(name) > 0 {
		fullname := filepath.Join(lockFilePath, name)
		// open処理移動
		var fd int
		var err error
		if _, err = os.Stat(fullname); err != nil {
			fd, err = syscall.Open(fullname, syscall.O_CREAT|syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
		} else {
			fd, err = syscall.Open(fullname, syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
		}
		if err != nil {
			return nil, err
		}
		return &LockHandle{fd, false}, nil
	} else {
		return &LockHandle{0, false}, errors.New("Invalid lockfile name.")
	}
}

// ファイルを利用して、ロックを行います。
// 引数で指定したミリ秒まで待機します。0以下を指定した場合は、リトライしません。
// 他プロセスのロックが指定時間内に解けなかった場合は、ErrBusy を返します。
func (l *LockHandle) Lock(timeout_millisec int) error {
	if l.fd == 0 {
		return errors.New("Not initialize.")
	}
	err := l.tryLock()

	if err == nil {
		return nil

	} else { // Locked by other process.
		if timeout_millisec > 0 {
			st := time.Now()
			for {
				time.Sleep(1 * time.Millisecond)
				err = l.tryLock()
				if err == nil {
					return nil // ロック成功
				}
				if time.Since(st).Nanoseconds() > (int64(timeout_millisec) * 1000000) {
					fmt.Fprintf(os.Stderr, "Lock Timeout %v\n", err)
					break
				}
			}
		}
	}
	return ErrBusy
}

// ロック解除。
func (l *LockHandle) Unlock() error {
	if !l.isLock {
		return errors.New("It has not been locked yet.")
	}
	if err := syscall.Flock(l.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	l.isLock = false
	return nil
}

// ロックファイルの終了処理。InitLock()成功後は、必ず呼び出して下さい。
func (l *LockHandle) TermLock() error {
	if l.fd != 0 {
		if l.isLock {
			l.Unlock()
		}
		syscall.Close(l.fd)
		l.fd = 0
	}
	return nil
}

// 実際にロック処理を行う。
func (l *LockHandle) tryLock() error {
	if err := syscall.Flock(l.fd, syscall.LOCK_EX); err != nil {
		return err
	}
	l.isLock = true
	return nil
}
