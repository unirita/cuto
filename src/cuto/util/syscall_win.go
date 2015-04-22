// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"syscall"
)

// DLLハンドル
type cutoDLL struct {
	dll *syscall.DLL
}

var (
	kernel32_dll = loadDLL("kernel32.dll")

	procCreateMutexW        = kernel32_dll.findProc("CreateMutexW")
	procWaitForSingleObject = kernel32_dll.findProc("WaitForSingleObject")
	procReleaseMutex        = kernel32_dll.findProc("ReleaseMutex")
	procCloseHandle         = kernel32_dll.findProc("CloseHandle")
	procGetModuleFileNameW  = kernel32_dll.findProc("GetModuleFileNameW")
	// WindowsAPIの関数アドレスはここに追加してください。
)

func loadDLL(name string) *cutoDLL {
	dll, err := syscall.LoadDLL(name)
	if err != nil {
		panic(err)
	}
	cutoDll := new(cutoDLL)
	cutoDll.dll = dll
	return cutoDll
}

func (c *cutoDLL) findProc(name string) *syscall.Proc {
	proc, err := c.dll.FindProc(name)
	if err != nil {
		panic(err)
	}
	return proc
}
