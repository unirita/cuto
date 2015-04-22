package remote

import (
	"net"
	"testing"
	"time"

	"cuto/servant/config"
)

type testConn struct {
	Written string
}

func init() {
	config.ReadConfig()
}

func (c *testConn) Close() error                       { return nil }
func (c *testConn) LocalAddr() net.Addr                { return nil }
func (c *testConn) RemoteAddr() net.Addr               { return nil }
func (c *testConn) SetDeadline(t time.Time) error      { return nil }
func (c *testConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *testConn) SetWriteDeadline(t time.Time) error { return nil }

func (c *testConn) Read(b []byte) (int, error) {
	msgBytes := []byte(testReqMsg)
	for i, c := range msgBytes {
		b[i] = c
	}

	return len(testReqMsg), nil
}

const testReqMsg = `{"type":"request","id":1234,"path":"C:\\work\\test.bat","param":"test","workspace": "C:\\work"}`

func (c *testConn) Write(b []byte) (int, error) {
	c.Written = string(b)
	return len(c.Written), nil
}

func TestStart_ポート番号に定義外の値を渡すとエラーが発生する(t *testing.T) {
	_, err := StartReceive(config.Servant.Sys.BindAddress, 65536, config.Servant.Job.MultiProc)

	if err == nil {
		t.Error("エラーが発生しませんでした。")
	}
}

func TestReceiveMessage_セッションキューにセッションを追加できる(t *testing.T) {
	conn := new(testConn)
	sq := make(chan *Session)

	go receiveMessage(conn, sq)

	select {
	case session := <-sq:
		if session.Conn == nil {
			t.Error("セッションにコネクションオブジェクトがセットされていません。")
		}

		if session.Body != testReqMsg {
			t.Error("セッションにセットされたメッセージが間違っています。")
			t.Logf("想定値: %s", testReqMsg)
			t.Logf("実績値: %s", session.Body)
		}
	case <-time.After(time.Second * 3):
		t.Error("3秒待ちましたが、セッションがキューに挿入されませんでした")
	}
}
