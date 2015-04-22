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

	isLock, err := hLock.Lock(2)
	if err != nil {
		t.Errorf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
	}
	if !isLock {
		t.Error("ロックに成功すべき所、失敗しました。")
	}
}

func TestUnlock_アンロックできる(t *testing.T) {
	hLock, err := InitMutex("testUnlock")
	if hLock == nil {
		t.Fatal(err)
		return
	}
	defer hLock.TermMutex()

	isLock, err := hLock.Lock(2)
	if err != nil {
		t.Fatalf("ロックに成功すべきところ、失敗しました。 - %v\n", err)
	}
	if !isLock {
		t.Fatal("ロックに成功すべきところ、失敗しました。")
	}
	if err = hLock.Unlock(); err != nil {
		t.Errorf("アンロックに成功すべきところ、失敗しました。 - %s", err)
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

var msg string

func parallel(wg *sync.WaitGroup) {

	hLock, err := InitMutex("testGo")
	if hLock == nil {
		panic(err)
		return
	}
	defer hLock.TermMutex()

	if ok, err := hLock.Lock(1000); !ok {
		panic(err)
	}
	defer hLock.Unlock()

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
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	wg := new(sync.WaitGroup)

//	for i := 0; i < 3; i++ {
//		wg.Add(1)
//		go parallel(wg)
//	}
//	wg.Wait()
//@todo
//	if msg != "abcdabcdabcd" {
//		t.Errorf("同期処理が正常に動作していないため、不正なメッセージが生成されました。[%v]\n", msg)
//	}
//}
