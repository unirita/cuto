// Copyright 2015 unirita Inc.
// Created 2015/06/03 shanxia

package util

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"io/ioutil"

	"path/filepath"
)

type LockHandle string

var (
	ErrBusy       = errors.New("Locked by other process.")
	errInvalidPid = errors.New("Lockfile contains invalid pid.")
	errDeadOwner  = errors.New("Lockfile contains pid of process not existent on this system.")

	lockFilePath = fmt.Sprintf("%s%c%s%c", GetRootPath(), os.PathSeparator, "temp", os.PathSeparator)

	onceLockFile = getOnceLock()
)

func getOnceLock() string {
	return filepath.Join(GetRootPath(), "temp", "once.lock")
}

// ファイルを利用した同期処理機能の初期化関数。
// ファイル作成が可能なファイル名を指定します。
func InitLock(name string) (*LockHandle, error) {
	if len(name) > 0 {
		lh := LockHandle(fmt.Sprintf("%s%s", lockFilePath, name))
		return &lh, nil
	} else {
		lh := LockHandle("")
		return &lh, errors.New("Invalid lockfile name.")
	}
}

// ファイルを利用して、ロックを行います。
// 引数で指定したミリ秒まで待機します。0以下を指定した場合は、1度だけロックに挑戦します。
// 他プロセスのロックが指定時間内に解けなかった場合は、ErrBusy を返します。
func (fl *LockHandle) Lock(timeout_millisec int) error {
	err := fl.tryLock(true)

	if err == ErrBusy { // Locked by other process.
		if timeout_millisec > 0 {
			// 既に他プロセスがロックしているので、指定時間リトライする。
			st := time.Now()
			for {
				time.Sleep(1 * time.Millisecond)
				err = fl.tryLock(true)
				if err == nil {
					return nil
				} else if err != ErrBusy {
					return err
				}
				d := time.Since(st)
				if d.Nanoseconds() > (int64(timeout_millisec) * 1000000) {
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
	name := string(*fl)

	content, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	var pid int
	_, err = fmt.Sscanln(string(content), &pid)
	if err != nil || pid < 1 {
		return errInvalidPid
	}
	if os.Getpid() != pid {
		return errors.New("Not locked.")
	}
	return os.Remove(name)
}

// ロックファイルの終了処理。（現在は何も行わない）
func (fl *LockHandle) TermLock() error {
	*fl = LockHandle("")
	return nil
}

// ロック処理を行う
func (fl *LockHandle) tryLock(execOnceLock bool) error {

	name := string(*fl)
	if len(name) < 1 {
		return errors.New("Not initialize.")
	}
	//@DEBUG start
	if execOnceLock {
		fmt.Fprintln(os.Stderr, "Before stat")
		if _, err := os.Stat(onceLockFile); err != nil {
			panic("Nothing " + onceLockFile)
		}
		fmt.Fprintln(os.Stderr, "Before open")
		fd, err := syscall.Open(onceLockFile, syscall.O_RDONLY|syscall.O_CLOEXEC, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(os.Stderr, "Before lock")
		if err := syscall.Flock(fd, syscall.LOCK_EX); err != nil {
			panic(err)
		}
		defer func() {
			fmt.Fprintln(os.Stderr, "Before unlock")
			syscall.Flock(fd, syscall.LOCK_UN)
			syscall.Close(fd)
		}()
	}
	//@DEBUG end

	tmpfile, err := ioutil.TempFile(filepath.Dir(name), "cuto_")
	if err != nil {
		return err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	_, err = tmpfile.WriteString(fmt.Sprintf("%d\n", os.Getpid()))
	if err != nil {
		return err
	}

	os.Link(tmpfile.Name(), name)

	fiTmp, err := os.Lstat(tmpfile.Name())
	if err != nil {
		return err
	}
	fiLock, err := os.Lstat(name)
	if err != nil {
		return err
	}
	if os.SameFile(fiTmp, fiLock) {
		// Successful.
		return nil
	}
	exist, err := fl.existOwnerProcess()
	if exist {
		return ErrBusy
	}
	// プロセスが存在しない理由によっては、ロックを開始します
	switch err {
	case errDeadOwner, errInvalidPid:
		// ゴミ情報を削除して、ロックに再挑戦します。
		err = os.Remove(name)
		if err != nil {
			return err
		}
		return fl.tryLock(false)
	default:
		return err
	}
	panic("Not reached.")
}

// ロックファイルには、ロックしたプロセスのIDのみを記録しているので、
// そのプロセスがまだ存在するか確認し、存在しない場合はerrorで理由を返します。
func (fl *LockHandle) existOwnerProcess() (bool, error) {
	name := string(*fl)

	content, err := ioutil.ReadFile(name)
	if err != nil {
		return false, err
	}

	var pid int
	_, err = fmt.Sscanln(string(content), &pid)
	if err != nil || pid < 0 {
		return false, errInvalidPid
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		if isAccessDenied(err) {
			// アクセスエラー時は、プロセスが存在する。
			return true, nil
		}
		return false, errDeadOwner
	}
	// 取得したプロセスに、シグナルを送ってみて、成功すれば生きていると判断します。
	err = p.Signal(os.Signal(syscall.Signal(0)))
	if err == nil {
		return true, nil
	}
	// zombie?
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, errDeadOwner
	}
	switch errno {
	case syscall.ESRCH:
		return false, errDeadOwner
	case syscall.EPERM:
		fmt.Fprintf(os.Stderr, "Operation not permitted from lock process[%d].\n", p.Pid)
		return true, nil
	default:
		return false, errors.New("Unknown error.")
	}
	panic("Not reached.")
}

// エラー情報が、権限エラーによる物か確認します。
func isAccessDenied(err error) bool {
	return false
}
