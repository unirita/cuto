package testutil

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestNewListenerStub_Listener型のスタブオブジェクトを生成できる(t *testing.T) {
	l := NewListenerStub()
	_, ok := interface{}(l).(net.Listener)
	if !ok {
		t.Fatalf("生成されたオブジェクトがnet.Listener型ではない。")
	}
}

func TestNewConnStub_Conn型のスタブオブジェクトを生成できる(t *testing.T) {
	c := NewConnStub()
	_, ok := interface{}(c).(net.Conn)
	if !ok {
		t.Fatalf("生成されたオブジェクトがnet.Conn型ではない。")
	}
}

func TestNewAddrStub_Addr型のスタブオブジェクトを生成できる(t *testing.T) {
	a := NewAddrStub("127.0.0.1:12345")
	_, ok := interface{}(a).(net.Addr)
	if !ok {
		t.Fatalf("生成されたオブジェクトがnet.Addr型ではない。")
	}
}

func TestListenerAccept_ConnStubオブジェクトを取得できる(t *testing.T) {
	l := NewListenerStub()
	c, err := l.Accept()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if c == nil {
		t.Errorf("オブジェクトが取得できていない。")
	}
}

func TestListenerAccept_エラーを発生させられる(t *testing.T) {
	l := NewListenerStub()
	l.AcceptErr = errors.New("error")
	_, err := l.Accept()
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestListenerClose_正常クローズできる(t *testing.T) {
	l := NewListenerStub()
	err := l.Close()
	if err != nil {
		t.Errorf("想定外のエラーが発生した: %s", err)
	}
	if !l.IsClosed {
		t.Errorf("クローズされていないことになっている。")
	}
}

func TestListenerClose_エラーを発生させられる(t *testing.T) {
	l := NewListenerStub()
	l.CloseErr = errors.New("error")
	err := l.Close()
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestListenerAddr_アドレス情報を取得できる(t *testing.T) {
	l := NewListenerStub()
	a := l.Addr()
	if a == nil {
		t.Fatalf("アドレス情報を取得できなかった。")
	}
}

func TestConnRead_読み込みができる(t *testing.T) {
	c := NewConnStub()
	c.ReadStr = "readtest"
	buf := make([]byte, 20)
	n, err := c.Read(buf)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	readStr := string(buf[:n])
	if readStr != "readtest" {
		t.Errorf("読み込んだ文字列[%s]が想定値と違う。", readStr)
	}
}

func TestConnRead_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.ReadStr = "readtest"
	c.ReadErr = errors.New("error")
	buf := make([]byte, 20)
	_, err := c.Read(buf)
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestConnWrite_書き込みができる(t *testing.T) {
	c := NewConnStub()
	buf := []byte("writetest")
	n, err := c.Write(buf)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if c.WriteStr != "writetest" {
		t.Errorf("書き込まれた文字列[%s]が想定値と違う。", c.WriteStr)
	}
	if n != len("writetest") {
		t.Errorf("書き込まれた長さ[%d]が想定値と違う。", n)
	}
}

func TestConnWrite_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.WriteErr = errors.New("error")
	buf := []byte("writetest")
	_, err := c.Write(buf)
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestConnClose_正常クローズできる(t *testing.T) {
	c := NewConnStub()
	err := c.Close()
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if !c.IsClosed {
		t.Errorf("クローズされていないことになっている。")
	}
}

func TestConnClose_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.CloseErr = errors.New("error")
	err := c.Close()
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestConnAddr_ローカルアドレス情報を取得できる(t *testing.T) {
	c := NewConnStub()
	a := c.LocalAddr()
	if a == nil {
		t.Fatalf("アドレス情報を取得できなかった。")
	}
}

func TestConnAddr_リモートアドレス情報を取得できる(t *testing.T) {
	c := NewConnStub()
	a := c.RemoteAddr()
	if a == nil {
		t.Fatalf("アドレス情報を取得できなかった。")
	}
}

func TestConnSetDeadline_正常な送受信待ち期限セット処理ができる(t *testing.T) {
	c := NewConnStub()
	err := c.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
}

func TestSetDeadline_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.SetDeadlineErr = errors.New("error")
	err := c.SetDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestConnSetReadDeadline_正常な受信待ち期限セット処理ができる(t *testing.T) {
	c := NewConnStub()
	err := c.SetReadDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
}

func TestSetReadDeadline_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.SetReadDeadlineErr = errors.New("error")
	err := c.SetReadDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestConnSetWriteDeadline_正常な送受信待ち期限セット処理ができる(t *testing.T) {
	c := NewConnStub()
	err := c.SetWriteDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
}

func TestSetWriteDeadline_エラーを発生させられる(t *testing.T) {
	c := NewConnStub()
	c.SetWriteDeadlineErr = errors.New("error")
	err := c.SetWriteDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Fatalf("エラーが発生していない。")
	}
}

func TestAddrNetwork_ネットワークの名前を取得できる(t *testing.T) {
	addrStr := "127.0.0.1:12345"
	a := NewAddrStub(addrStr)
	if a.Network() != "tcp" {
		t.Errorf("取得した値[%s]が想定と違っている。", a.Network())
	}
}

func TestAddrString_アドレスの文字列表現を取得できる(t *testing.T) {
	addrStr := "127.0.0.1:12345"
	a := NewAddrStub(addrStr)
	if a.String() != addrStr {
		t.Errorf("取得した値[%s]が想定と違っている。", a.String())
	}
}
