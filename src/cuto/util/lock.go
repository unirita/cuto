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
type MutexHandle struct {
	handle uintptr
	isLock bool
}

const (
	wAIT_OBJECT_0  int = 0
	wAIT_ABANDONED int = 128
	wAIT_TIMEOUT   int = 258
)

// プロセス間で共通に使用する名前を指定する。
func InitMutex(name string) (*MutexHandle, error) {
	mutexName := fmt.Sprintf("Global\\%s", name)
	hMutex, _, err := procCreateMutexW.Call(
		0, 0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(mutexName))))
	if hMutex == 0 {
		fmt.Fprintf(os.Stderr, "Failed InitMutexW() err = %v", err)
		return nil, err
	}
	return &MutexHandle{hMutex, false}, nil
}

// ロックを開始する。
// 引数でタイムアウト時間（ミリ秒）を指定する。
func (m *MutexHandle) Lock(timeout_milisec int) (bool, error) {
	r1, _, err := procWaitForSingleObject.Call(m.handle, uintptr(timeout_milisec))
	if int(r1) == wAIT_OBJECT_0 || int(r1) == wAIT_ABANDONED {
		// Lock成功
		m.isLock = true
		return true, nil
	} else if int(r1) == wAIT_TIMEOUT {
		msg := fmt.Sprintf("Lock Timeout. Detail( %v )", err)
		fmt.Fprintf(os.Stderr, "%v\n", msg)
		return false, errors.New(msg)
	}
	return false, fmt.Errorf("Lock Unknown Error. Detail( %v )", err)
}

// ロック中であれば、解除する。
func (m *MutexHandle) Unlock() error {
	if m.isLock {
		r1, _, err := procReleaseMutex.Call(m.handle)
		if int(r1) == 0 { // 失敗
			return fmt.Errorf("Unlock Error. Detail( %v )", err)
		}
		m.isLock = false
		return nil
	}
	return nil
}

// 自プロセスがロックしているかの確認
func (m *MutexHandle) IsLock() bool {
	isLock := false
	r1, _, _ := procWaitForSingleObject.Call(m.handle, 0)
	if int(r1) == wAIT_OBJECT_0 || int(r1) == wAIT_ABANDONED {
		isLock = m.isLock
		procReleaseMutex.Call(m.handle)
		return isLock
	}
	return isLock
}

// InitMutexで確保したミューテックスオブジェクトを破棄する。
func (m *MutexHandle) TermMutex() error {
	procCloseHandle.Call(m.handle)
	return nil
}
