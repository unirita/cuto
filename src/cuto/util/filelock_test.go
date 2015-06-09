// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"cuto/testutil"
)

var (
	lockFile = createLockName()
)

func createLockName() string {
	var f string
	d := os.Getenv("GOPATH")
	if runtime.GOOS == "windows" {
		f = fmt.Sprintf("%s\\test\\cuto\\util\\test.lock", d)
	} else {
		f = fmt.Sprintf("%s/test/cuto/util/test.lock", d)
	}
	if _, err := os.Stat(f); err == nil {
		os.Remove(f)
	}
	return f
}

func TestNewFileLock_初期化できる(t *testing.T) {
	_, err := NewFileLock(lockFile)
	if err != nil {
		t.Errorf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
}

func TestNewFileLock_初期化に失敗(t *testing.T) {
	_, err := NewFileLock("test.lock")
	if err == nil {
		t.Error("同期処理の初期化に成功しました。")
	}
}

func TestLock_初期化せずにロック(t *testing.T) {
	l, err := NewFileLock("test.lock")
	if err != nil {
		err := l.Lock(0)
		if err == nil {
			t.Error("初期化せずにロックしたが、成功した。")
		}
	} else {
		t.Error("同期処理の初期化に失敗しました。 - %v", err)
	}
}

func TestLock_ロックする(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	err = l.Lock(0)
	if err != nil {
		t.Errorf("ロックに失敗しました。 - %v", err)
	}
	defer l.Unlock()

	if _, err = os.Stat(lockFile); err != nil {
		t.Errorf("ロックファイルが存在しない。 - %v", err)
	}
}

func TestLock_上位の権限プロセスがロック中のためにロックに失敗する(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	if runtime.GOOS == "windows" {
		file.WriteString("4") // system
	} else {
		file.WriteString("1") // init.d or launch.d or system
	}
	file.Close()

	c := testutil.NewStderrCapturer()
	c.Start()
	defer c.Stop()

	err = l.Lock(0)
	if err == nil {
		t.Error("他プロセスのロック中に、ロック成功しました。")
		defer l.Unlock()
	} else if err != ErrBusy {
		t.Errorf("予期しないエラーが返りました。 - %v", err)
	}
}

func TestLock_他プロセスがロック中なのでロックに失敗する(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	file.WriteString("568") // init.d or launch.d or system
	file.Close()

	c := testutil.NewStderrCapturer()
	c.Start()
	defer c.Stop()

	err = l.Lock(100)
	if err == nil {
		t.Error("他プロセスのロック中に、ロック成功しました。")
		defer l.Unlock()
	} else if err != ErrBusy {
		t.Errorf("予期しないエラーが返りました。 - %v", err)
	}
}

func TestLock_ロックファイルが残っている状態でロックに成功する(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	file.WriteString("99999") // ありえないプロセス
	file.Close()

	c := testutil.NewStderrCapturer()
	c.Start()
	defer c.Stop()

	err = l.Lock(1)
	if err != nil && err != ErrBusy {
		t.Errorf("予期しないエラーが発生しました - %v", err)
	}
	if err == nil {
		l.Unlock()
	}
}

func TestLock_不正なロックファイルが残っている状態でロックに成功する(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	file.WriteString("A") // ありえないプロセス
	file.Close()

	c := testutil.NewStderrCapturer()
	c.Start()
	defer c.Stop()

	err = l.Lock(1)
	if err != nil && err != ErrBusy {
		t.Errorf("予期しないエラーが発生しました - %v", err)
	}
	if err == nil {
		l.Unlock()
	}
}

func TestUnlock_初期化前にアンロックする(t *testing.T) {
	l, err := NewFileLock("test.lock")
	if err != nil {
		err := l.Unlock()
		if err == nil {
			t.Error("初期化せずにアンロックしたが、成功した。")
		}
	} else {
		t.Error("同期処理の初期化に失敗しました。 - %v", err)
	}
}

func TestUnlock_ロック前にアンロックする(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	err = l.Unlock()
	if err == nil {
		t.Error("アンロックに失敗すべきところ、成功しました。")
	}
}

func TestUnlock_他プロセスがロック中なのでアンロックに失敗(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	file.WriteString("99999") // ありえないプロセス
	file.Close()

	err = l.Unlock()
	if err == nil {
		t.Error("アンロックに失敗すべきところ、成功しました。")
	} else if err.Error() != "Not locked." {
		t.Errorf("想定外のエラーが返りました。 - %v", err)
	}
}

func TestUnlock_不正なロック中でアンロックに失敗(t *testing.T) {
	l, err := NewFileLock(lockFile)
	if err != nil {
		t.Fatalf("同期処理の初期化に失敗しました。 - %s", err.Error())
	}
	// 他プロセスによってロック
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("ロックファイルの作成失敗。 - %v", err)
	}
	defer os.Remove(lockFile)
	file.WriteString("S") // ありえないプロセス
	file.Close()

	err = l.Unlock()
	if err == nil {
		t.Error("アンロックに失敗すべきところ、成功しました。")
	} else if err != errInvalidPid {
		t.Errorf("想定外のエラーが返りました。 - %v", err)
	}
}
