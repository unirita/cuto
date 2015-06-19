// Copyright 2015 unirita Inc.
// Created 2015/06/03 shanxia

package util

import (
	"errors"
	"os"
	"syscall"
	"time"

	"path/filepath"
)

type LockHandle struct {
	name string //lock ファイル名
	fd   int    //lockファイルディスクリプタ
}

var (
	ErrBusy = errors.New("Locked by other process.")

	lockFilePath = filepath.Join(GetRootPath(), "temp")
)

// ファイルを利用した同期処理機能の初期化関数。
// ファイル作成が可能なファイル名を指定します。
func InitLock(name string) (*LockHandle, error) {
	if len(name) > 0 {
		return &LockHandle{filepath.Join(lockFilePath, name), 0}, nil
	} else {
		return &LockHandle{"", 0}, errors.New("Invalid lockfile name.")
	}
}

// ファイルを利用して、ロックを行います。
// 引数で指定したミリ秒まで待機します。0以下を指定した場合は、1度だけロックに挑戦します。
// 他プロセスのロックが指定時間内に解けなかった場合は、ErrBusy を返します。
func (fl *LockHandle) Lock(timeout_millisec int) error {
	err := fl.tryLock()

	if err == ErrBusy { // Locked by other process.
		if timeout_millisec > 0 {
			// 既に他プロセスがロックしているので、指定時間リトライする。
			st := time.Now()
			for {
				time.Sleep(1 * time.Millisecond)
				err = fl.tryLock()
				if err == nil {
					return nil
				} else if err != ErrBusy {
					return err
				}
				if time.Since(st).Nanoseconds() > (int64(timeout_millisec) * 1000000) {
					return ErrBusy
				}
			}
		}
		return ErrBusy
	} else {
		return err
	}
	panic("Not reached.")
}

// ロック解除。
func (fl *LockHandle) Unlock() error {
	if fl.fd == 0 {
		return errors.New("It has not been locked yet.")
	}
	defer syscall.Close(fl.fd)
	if err := syscall.Flock(fl.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return nil
}

// ロックファイルの終了処理。（現在は何も行わない）
func (fl *LockHandle) TermLock() error {
	fl.Unlock()
	return nil
}

// ロック処理を行う
func (fl *LockHandle) tryLock() error {
	if len(fl.name) < 1 {
		return errors.New("Not initialize.")
	}
	var err error
	if _, err := os.Stat(fl.name); err != nil {
		fl.fd, err = syscall.Open(fl.name, syscall.O_CREAT|syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
	} else {
		fl.fd, err = syscall.Open(fl.name, syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
	}
	if err != nil {
		return err
	}
	if err := syscall.Flock(fl.fd, syscall.LOCK_EX); err != nil {
		return ErrBusy
	}
	return nil
}
