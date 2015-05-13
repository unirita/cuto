// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	_ "runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestInitMutex_初期化できる(t *testing.T) {
	hLock, err := InitMutex("testInit")
	if err != nil {
		t.Errorf("初期化に失敗しました。 - %v\n", err)
	}
	if hLock == nil {
		t.Error("排他オブジェクトがnilを返しました。")
	}
	defer hLock.TermMutex()
}

func TestTerm_破棄できる(t *testing.T) {
	hLock, err := InitMutex("testTerm")
	if hLock == nil {
		t.Fatal(err)
		return
	}

	if err = hLock.TermMutex(); err != nil {
		t.Fatal(err)
	}
}

func TestLock_ロックできる(t *testing.T) {
	hLock, err := InitMutex("testLock")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermMutex()

	if hLock.IsLock() {
		t.Error("ロック前にもかかわらず、trueが返りました。")
	}
	ok, err := hLock.Lock(2)
	if err != nil {
		t.Errorf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
	}
	defer hLock.Unlock()
	if !ok {
		t.Error("ロックに成功すべき所、失敗しました。")
	}
	if !hLock.IsLock() {
		t.Error("ロック中にもかかわらず、falseが返りました。")
	}
}

func TestUnlock_アンロックできる(t *testing.T) {
	hLock, err := InitMutex("testUnlock")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermMutex()

	ok, err := hLock.Lock(2)
	if err != nil {
		t.Fatalf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
	}
	if !ok {
		t.Fatal("ロックに成功すべきところ、失敗しました。")
	}
	if !hLock.IsLock() {
		t.Error("ロック中にもかかわらず、falseが返りました。")
	}
	if err = hLock.Unlock(); err != nil {
		t.Errorf("アンロックに成功すべきところ、失敗しました。 - %s", err)
	}
	if hLock.IsLock() {
		t.Error("アンロック後にもかかわらず、trueが返りました。")
	}
}

func TestUnlock_ロックしていないがアンロックする(t *testing.T) {
	hLock, err := InitMutex("testUnlock2")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermMutex()

	if err = hLock.Unlock(); err != nil {
		t.Errorf("アンロックに成功すべきところ、失敗しました。 - %s", err)
	}
}

func TestUnlock_ロックしていないが無理矢理アンロック処理をする(t *testing.T) {
	hLock, err := InitMutex("testUnlock2")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermMutex()

	hLock.isLock = true
	if err = hLock.Unlock(); err == nil {
		t.Errorf("アンロックに失敗すべきところ、成功しました。")
	} else if !strings.Contains(err.Error(), "Attempt to release mutex not owned by caller") {
		t.Errorf("予想外のエラーが返りました。 - %v", err)
	}
}

var msg string

func parallel(wg *sync.WaitGroup, m *sync.Mutex, hLock *MutexHandle) {

	m.Lock()

	if ok, err := hLock.Lock(10000); !ok {
		panic(err)
	}
	defer hLock.Unlock()

	m.Unlock()

	msg += "a"
	time.Sleep(100 * time.Millisecond)
	msg += "b"
	time.Sleep(100 * time.Millisecond)
	msg += "c"
	time.Sleep(100 * time.Millisecond)
	msg += "d"

	wg.Done()
}

//func TestLock_goルーチンを利用した排他確認(t *testing.T) {
//	hLock, err := InitMutex("testGo")
//	if hLock == nil {
//		panic(err)
//		return
//	}
//	defer hLock.TermMutex()

//	runtime.GOMAXPROCS(runtime.NumCPU())
//	wg := new(sync.WaitGroup)

//	lock, err := InitMutex("testGo")
//	if err != nil {
//		panic(err)
//		return
//	}
//	defer lock.TermMutex()

//	var m sync.Mutex
//	for i := 0; i < 3; i++ {
//		wg.Add(1)
//		go parallel(wg, &m, hLock)
//	}
//	wg.Wait()

//	if msg != "abcdabcdabcd" {
//		t.Errorf("同期処理が正常に動作していないため、不正なメッセージが生成されました。[%v]\n", msg)
//	}
//}
