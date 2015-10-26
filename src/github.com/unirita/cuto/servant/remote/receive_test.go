package remote

import (
	"errors"
	"testing"

	"github.com/unirita/cuto/servant/config"
	"github.com/unirita/cuto/testutil"
)

func init() {
	config.ReadConfig("")
}

func TestStartReceive_ポート番号に定義外の値を渡すとエラーが発生する(t *testing.T) {
	_, err := StartReceive(config.Servant.Sys.BindAddress, 65536, config.Servant.Job.MultiProc)

	if err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestReceiveLoopProcess_セッションキューにセッションを追加できる(t *testing.T) {
	listener := testutil.NewListenerStub()
	sq := make(chan *Session, 1)
	err := receiveLoopProcess(listener, sq)
	defer close(sq)

	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	session := <-sq
	if session == nil {
		t.Error("セッションがキューに挿入されていない。")
	}
}

func TestReceiveLoopProcess_Acceptに失敗したらエラー(t *testing.T) {
	listener := testutil.NewListenerStub()
	listener.AcceptErr = errors.New("testerror")
	sq := make(chan *Session, 1)
	err := receiveLoopProcess(listener, sq)
	close(sq)

	if err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestReceiveMessage_セッションキューにセッションを追加できる(t *testing.T) {
	reqMsg := `{"type":"request","id":1234,"path":"C:\\work\\test.bat","param":"test","workspace": "C:\\work"}`
	conn := testutil.NewConnStub()
	conn.ReadStr = reqMsg + "\n"
	sq := make(chan *Session, 1)
	err := receiveMessage(conn, sq)
	close(sq)

	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}

	session := <-sq
	if session == nil {
		t.Error("セッションがキューに挿入されていない。")
	}

	if session.Conn == nil {
		t.Error("セッションにコネクションオブジェクトがセットされていない。")
	}

	if session.Body != reqMsg {
		t.Error("セッションにセットされたメッセージが間違っている。")
		t.Logf("想定値: %s", reqMsg)
		t.Logf("実績値: %s", session.Body)
	}
}

func TestReceiveMessage_待ち期限の設定に失敗したらエラー(t *testing.T) {
	conn := testutil.NewConnStub()
	conn.ReadStr = `{"type":"request","id":1234,"path":"C:\\work\\test.bat","param":"test","workspace": "C:\\work"}`
	conn.SetReadDeadlineErr = errors.New("testerror")
	sq := make(chan *Session, 1)
	err := receiveMessage(conn, sq)
	close(sq)

	if err == nil {
		t.Fatal("エラーが発生していない。")
	}

	session := <-sq
	if session != nil {
		t.Error("エラーが発生したにも関わらず、セッションがキューに挿入された。")
	}
}

func TestReceiveMessage_メッセージの読み込みに失敗したらエラー(t *testing.T) {
	conn := testutil.NewConnStub()
	conn.ReadStr = `{"type":"request","id":1234,"path":"C:\\work\\test.bat","param":"test","workspace": "C:\\work"}`
	conn.ReadErr = errors.New("testerror")
	sq := make(chan *Session, 1)
	err := receiveMessage(conn, sq)
	close(sq)

	if err == nil {
		t.Fatal("エラーが発生していない。")
	}

	session := <-sq
	if session != nil {
		t.Error("エラーが発生したにも関わらず、セッションがキューに挿入された。")
	}
}
