// Copyright 2015 unirita Inc.
// Created 2015/06/03 shanxia

package util

import (
	"errors"
	"fmt"
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
		fullname := filepath.Join(lockFilePath, name)
		// open処理移動
		var fd int
		var err error
		if _, err = os.Stat(fullname); err != nil {
			fmt.Printf("Create File %v\n", fullname)
			fd, err = syscall.Open(fullname, syscall.O_CREAT|syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
		} else {
			fmt.Printf("Open File %v\n", fullname)
			fd, err = syscall.Open(fullname, syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
		}
		if err != nil {
			return nil, err
		}
		return &LockHandle{filepath.Join(lockFilePath, name), fd}, nil
		//		return &LockHandle{filepath.Join(lockFilePath, name), 0}, nil
	} else {
		return &LockHandle{"", 0}, errors.New("Invalid lockfile name.")
	}
}

// ファイルを利用して、ロックを行います。
// 引数で指定したミリ秒まで待機します。0以下を指定した場合は、1度だけロックに挑戦します。
// 他プロセスのロックが指定時間内に解けなかった場合は、ErrBusy を返します。
func (l *LockHandle) Lock(timeout_millisec int) error {
	err := l.tryLock()

	if err == nil {
		return nil

	} else if err == ErrBusy { // Locked by other process.
		if timeout_millisec > 0 {
			st := time.Now()
			for {
				time.Sleep(1 * time.Millisecond)
				err = l.tryLock()
				if err == nil {
					return nil // ロック成功
				}
				if time.Since(st).Nanoseconds() > (int64(timeout_millisec) * 1000000) {
					break
				}
			}
		}
		syscall.Close(l.fd)
		l.fd = 0
		return ErrBusy
	}
	return err
}

// ロック解除。
func (l *LockHandle) Unlock() error {
	if l.fd == 0 {
		return errors.New("It has not been locked yet.")
	}
	defer func() {
		//		syscall.Close(l.fd)
		//		l.fd = 0
	}()
	if err := syscall.Flock(l.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	return nil
}

// ロックファイルの終了処理。（現在は何も行わない）
func (l *LockHandle) TermLock() error {
	syscall.Close(l.fd)
	l.fd = 0
	return nil
}

// 実際にロック処理を行う。
// 現在は、nilまたはErrBusyを返すと、ファイルを開いている状態という目印にもなっている。
func (l *LockHandle) tryLock() error {
	if len(l.name) == 0 {
		return errors.New("Not initialize.")
	}
	//	if l.fd == 0 {
	//		var err error
	//		if _, err = os.Stat(l.name); err != nil {
	//			l.fd, err = syscall.Open(l.name, syscall.O_CREAT|syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
	//		} else {
	//			l.fd, err = syscall.Open(l.name, syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
	//		}
	//		if err != nil {
	//			return err
	//		}
	//	}
	if err := syscall.Flock(l.fd, syscall.LOCK_EX); err != nil {
		return ErrBusy
	}
	return nil
}
