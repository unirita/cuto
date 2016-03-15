// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package util

import (
	_ "runtime"
	"sync"
	"testing"
	"time"
)

func TestInitMutex_初期化できる(t *testing.T) {
	hLock, err := InitLock("testInit")
	if err != nil {
		t.Errorf("初期化に失敗しました。 - %v\n", err)
	}
	if hLock == nil {
		t.Error("排他オブジェクトがnilを返しました。")
	}
	defer hLock.TermLock()
}

func TestTerm_破棄できる(t *testing.T) {
	hLock, err := InitLock("testTerm")
	if hLock == nil {
		t.Fatal(err)
		return
	}

	if err = hLock.TermLock(); err != nil {
		t.Fatal(err)
	}
}

func TestLock_ロックできる(t *testing.T) {
	hLock, err := InitLock("testLock")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermLock()

	if hLock.IsLock() {
		t.Error("ロック前にもかかわらず、trueが返りました。")
	}
	err = hLock.Lock(2)
	if err != nil {
		t.Errorf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
	}
	defer hLock.Unlock()
	if !hLock.IsLock() {
		t.Error("ロック中にもかかわらず、falseが返りました。")
	}
}

func TestUnlock_アンロックできる(t *testing.T) {
	hLock, err := InitLock("testUnlock")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermLock()

	err = hLock.Lock(2)
	if err != nil {
		t.Fatalf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
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
	hLock, err := InitLock("testUnlock2")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermLock()

	if err = hLock.Unlock(); err != nil {
		t.Errorf("アンロックに成功すべきところ、失敗しました。 - %s", err)
	}
}

var msg string

func parallel(wg *sync.WaitGroup, m *sync.Mutex, hLock *LockHandle) {

	m.Lock()

	if err := hLock.Lock(10000); err != nil {
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
