package remote

import (
	"testing"
	"time"

	"cuto/servant/config"
	"cuto/testutil"
)

func init() {
	config.ReadConfig()
}

func TestStart_ポート番号に定義外の値を渡すとエラーが発生する(t *testing.T) {
	_, err := StartReceive(config.Servant.Sys.BindAddress, 65536, config.Servant.Job.MultiProc)

	if err == nil {
		t.Error("エラーが発生しませんでした。")
	}
}

func TestReceiveMessage_セッションキューにセッションを追加できる(t *testing.T) {
	conn := testutil.NewConnStub()
	conn.ReadStr = `{"type":"request","id":1234,"path":"C:\\work\\test.bat","param":"test","workspace": "C:\\work"}`
	sq := make(chan *Session)

	go receiveMessage(conn, sq)

	select {
	case session := <-sq:
		if session.Conn == nil {
			t.Error("セッションにコネクションオブジェクトがセットされていません。")
		}

		if session.Body != conn.ReadStr {
			t.Error("セッションにセットされたメッセージが間違っています。")
			t.Logf("想定値: %s", conn.ReadStr)
			t.Logf("実績値: %s", session.Body)
		}
	case <-time.After(time.Second * 3):
		t.Error("3秒待ちましたが、セッションがキューに挿入されませんでした")
	}
}
