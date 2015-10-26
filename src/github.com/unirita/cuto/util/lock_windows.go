// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// ミューテックスハンドルを保持する。
type LockHandle struct {
	handle uintptr
	isLock bool
}

const (
	wAIT_OBJECT_0  int = 0
	wAIT_ABANDONED int = 128
	wAIT_TIMEOUT   int = 258
)

var ErrBusy = errors.New("Locked by other process.")

// プロセス間で共通に使用する名前を指定する。
func InitLock(name string) (*LockHandle, error) {
	mutexName := fmt.Sprintf("Global\\%s", name)
	hMutex, _, err := procCreateMutexW.Call(
		0, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(mutexName))))
	if hMutex == 0 {
		fmt.Fprintf(os.Stderr, "Failed InitMutexW() err = %v", err)
		return nil, err
	}
	return &LockHandle{hMutex, false}, nil
}

// ロックを開始する。
// 引数でタイムアウト時間（ミリ秒）を指定する。
func (l *LockHandle) Lock(timeout_milisec int) error {
	r1, _, err := procWaitForSingleObject.Call(l.handle, uintptr(timeout_milisec))
	if int(r1) == wAIT_OBJECT_0 || int(r1) == wAIT_ABANDONED {
		// Lock成功
		l.isLock = true
		return nil
	} else if int(r1) == wAIT_TIMEOUT {
		msg := fmt.Sprintf("Lock Timeout. Detail( %v )", err)
		fmt.Fprintf(os.Stderr, "%v\n", msg)
		return ErrBusy
	}
	return fmt.Errorf("Lock Unknown Error. Detail( %v )", err)
}

// ロック中であれば、解除する。
func (l *LockHandle) Unlock() error {
	if l.isLock {
		r1, _, err := procReleaseMutex.Call(l.handle)
		if int(r1) == 0 { // 失敗
			return fmt.Errorf("Unlock Error. Detail( %v )", err)
		}
		l.isLock = false
		return nil
	}
	return nil
}

// 自プロセスがロックしているかの確認
func (l *LockHandle) IsLock() bool {
	isLock := false
	r1, _, _ := procWaitForSingleObject.Call(l.handle, 0)
	if int(r1) == wAIT_OBJECT_0 || int(r1) == wAIT_ABANDONED {
		isLock = l.isLock
		procReleaseMutex.Call(l.handle)
		return isLock
	}
	return isLock
}

// InitMutexで確保したミューテックスオブジェクトを破棄する。
func (l *LockHandle) TermLock() error {
	procCloseHandle.Call(l.handle)
	return nil
}
