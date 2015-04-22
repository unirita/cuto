package testutil

import (
	"net"
	"time"
)

// net.Listenerのスタブ構造体
type ListenerStub struct {
	AcceptErr error
	CloseErr  error

	addr net.Addr
}

// net.Connのスタブ構造体
type ConnStub struct {
	ReadStr             string // Read関数で取得させる文字列
	WriteStr            string // Write関数で書き込まれた文字列
	ReadErr             error
	WriteErr            error
	CloseErr            error
	SetDeadlineErr      error
	SetReadDeadlineErr  error
	SetWriteDeadlineErr error

	localAddr  net.Addr
	remoteAddr net.Addr
}

// net.Addrのスタブ構造体
type AddrStub struct {
	nwk string
	str string
}

func NewListenerStub() *ListenerStub {
	listener := new(ListenerStub)
	listener.addr = NewAddrStub("127.0.0.1:12345")
	return listener
}

func NewConnStub() *ConnStub {
	conn := new(ConnStub)
	conn.localAddr = NewAddrStub("127.0.0.1:54321")
	conn.remoteAddr = NewAddrStub("127.0.0.1:12345")
	return conn
}

func NewAddrStub(str string) *AddrStub {
	addr := new(AddrStub)
	addr.nwk = "tcp"
	addr.str = str
	return addr
}

func (l *ListenerStub) Accept() (c net.Conn, err error) {
	if l.AcceptErr != nil {
		return nil, l.AcceptErr
	}
	return NewConnStub(), nil
}

func (l *ListenerStub) Close() error {
	return l.CloseErr
}

func (l *ListenerStub) Addr() net.Addr {
	return l.addr
}

func (c *ConnStub) Read(b []byte) (n int, err error) {
	if c.ReadErr != nil {
		return 0, c.ReadErr
	}
	copy(b, []byte(c.ReadStr))
	return len(c.ReadStr), nil
}

func (c *ConnStub) Write(b []byte) (n int, err error) {
	if c.WriteErr != nil {
		return 0, c.WriteErr
	}
	c.WriteStr = string(b)
	return len(c.WriteStr), nil
}

func (c *ConnStub) Close() error {
	return c.CloseErr
}

func (c *ConnStub) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *ConnStub) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *ConnStub) SetDeadline(t time.Time) error {
	return c.SetDeadlineErr
}

func (c *ConnStub) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadlineErr
}

func (c *ConnStub) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadlineErr
}

func (a *AddrStub) Network() string {
	return a.nwk
}

func (a *AddrStub) String() string {
	return a.str
}
