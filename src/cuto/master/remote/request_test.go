package remote

import (
	"fmt"
	"net"
	"testing"
	"time"

	"cuto/master/config"
	"cuto/message"
)

func setTestConfig() {
	config.Job.ConnectionTimeoutSec = 1
}

func runTestReceiver(t *testing.T, listener net.Listener, msq chan<- string, delay int) {
	conn, err := listener.Accept()
	if err != nil {
		t.Log(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	readLen, err := conn.Read(buf)
	if err != nil {
		t.Log(err)
	}

	readMsg := string(buf[:readLen])
	msq <- readMsg

	d := time.Duration(delay) * time.Second
	time.Sleep(d)

	_, err = conn.Write([]byte(`testresponse`))
	if err != nil {
		t.Log(err)
	}
}

func runTestReceiverWithHearbeat(t *testing.T, listener net.Listener, msq chan<- string, delay int) {
	conn, err := listener.Accept()
	if err != nil {
		t.Log(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	readLen, err := conn.Read(buf)
	if err != nil {
		t.Log(err)
	}

	readMsg := string(buf[:readLen])
	msq <- readMsg

	s := 500 * time.Millisecond
	for i := 0; i < delay*2; i++ {
		time.Sleep(s)
		_, err = conn.Write([]byte(message.HEARTBEAT))
		if err != nil {
			t.Log(err)
		}
	}

	_, err = conn.Write([]byte(`testresponse`))
	if err != nil {
		t.Log(err)
	}
}

func runTestReceiverWithStartTime(t *testing.T, listener net.Listener, msq chan<- string, delay int) {
	conn, err := listener.Accept()
	if err != nil {
		t.Log(err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 1024)
	readLen, err := conn.Read(buf)
	if err != nil {
		t.Log(err)
	}

	readMsg := string(buf[:readLen])
	msq <- readMsg

	d := time.Duration(delay) * time.Second
	time.Sleep(d)

	conn.Write([]byte(message.ST_HEADER + "20150401123456.789"))

	_, err = conn.Write([]byte(`testresponse`))
	if err != nil {
		t.Log(err)
	}
}

func TestSendMessage_メッセージを送信できる(t *testing.T) {
	setTestConfig()
	host := "localhost"
	port := 12345
	hostPort := fmt.Sprintf("%s:%d", host, port)

	listener, listenErr := net.Listen("tcp", hostPort)
	if listenErr != nil {
		t.Fatalf("テスト用のlistenに失敗しました: %s", listenErr)
	}

	defer listener.Close()

	msq := make(chan string, 10)
	go runTestReceiver(t, listener, msq, 0)

	stCh := make(chan string, 1)
	defer close(stCh)
	resMsg, err := SendRequest(host, port, `testrequest`, stCh)
	if err != nil {
		t.Fatalf("エラーが発生しました: %s", err)
	}

	if message := <-msq; message != `testrequest` {
		t.Errorf("リスナに届いたメッセージが間違っています: %s", message)
	}

	if resMsg != `testresponse` {
		t.Errorf("リスナからのレスポンスメッセージが間違っています: %s", resMsg)
	}
}

func TestSendMessage_一定時間応答がない場合はタイムアウトする(t *testing.T) {
	setTestConfig()
	host := "localhost"
	port := 12345
	hostPort := fmt.Sprintf("%s:%d", host, port)

	listener, listenErr := net.Listen("tcp", hostPort)
	if listenErr != nil {
		t.Fatalf("テスト用のlistenに失敗しました: %s", listenErr)
	}

	defer listener.Close()

	msq := make(chan string, 10)
	go runTestReceiver(t, listener, msq, 2)

	stCh := make(chan string, 1)
	defer close(stCh)
	_, err := SendRequest(host, port, `testrequest`, stCh)
	if err == nil {
		t.Fatalf("タイムアウトが発生しない。")
	}
}

func TestSendMessage_ハートビートが返される場合はタイムアウトしない(t *testing.T) {
	setTestConfig()
	host := "localhost"
	port := 12345
	hostPort := fmt.Sprintf("%s:%d", host, port)

	listener, listenErr := net.Listen("tcp", hostPort)
	if listenErr != nil {
		t.Fatalf("テスト用のlistenに失敗しました: %s", listenErr)
	}

	defer listener.Close()

	msq := make(chan string, 10)
	go runTestReceiverWithHearbeat(t, listener, msq, 2)

	stCh := make(chan string, 1)
	defer close(stCh)
	resMsg, err := SendRequest(host, port, `testrequest`, stCh)
	if err != nil {
		t.Fatalf("エラーが発生しました: %s", err)
	}

	if message := <-msq; message != `testrequest` {
		t.Errorf("リスナに届いたメッセージが間違っています: %s", message)
	}

	if resMsg != `testresponse` {
		t.Errorf("リスナからのレスポンスメッセージが間違っています: %s", resMsg)
	}
}

func TestSendMessage_スタート時刻をチャンネルから取得できる(t *testing.T) {
	setTestConfig()
	host := "localhost"
	port := 12345
	hostPort := fmt.Sprintf("%s:%d", host, port)

	listener, listenErr := net.Listen("tcp", hostPort)
	if listenErr != nil {
		t.Fatalf("テスト用のlistenに失敗しました: %s", listenErr)
	}

	defer listener.Close()

	msq := make(chan string, 10)
	go runTestReceiverWithStartTime(t, listener, msq, 0)

	stCh := make(chan string, 1)
	defer close(stCh)
	resMsg, err := SendRequest(host, port, `testrequest`, stCh)
	if err != nil {
		t.Fatalf("エラーが発生しました: %s", err)
	}

	select {
	case st := <-stCh:
		if st != "20150401123456.789" {
			t.Errorf("取得したスタート時刻[%s]が間違っている。", st)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("十分な時間待ったが、スタート時刻が取得できなかった。")
	}

	if message := <-msq; message != `testrequest` {
		t.Errorf("リスナに届いたメッセージが間違っています: %s", message)
	}

	if resMsg != `testresponse` {
		t.Errorf("リスナからのレスポンスメッセージが間違っています: %s", resMsg)
	}
}
