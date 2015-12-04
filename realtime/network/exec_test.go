package network

import (
	"regexp"
	"testing"
	"time"
)

func TestNewCommand_WithEmptyRealtimeName(t *testing.T) {
	cmd := NewCommand("")

	pattern := `^realtime_\d{14}$`
	matcher := regexp.MustCompile(pattern)
	if !matcher.MatchString(cmd.GetNetworkName()) {
		t.Errorf("Unexpected Jobnet name format.")
		t.Logf("Jobnet name: %s", cmd.GetNetworkName())
	}
}

func TestNewCommand_WithRealtimeName(t *testing.T) {
	cmd := NewCommand("test")

	pattern := `^realtime_test_\d{14}$`
	matcher := regexp.MustCompile(pattern)
	if !matcher.MatchString(cmd.GetNetworkName()) {
		t.Errorf("Unexpected Jobnet name format.")
		t.Logf("Jobnet name: %s", cmd.GetNetworkName())
	}
}

func TestWaitID(t *testing.T) {
	cmd := new(Command)
	lineCh := make(chan string, 10)
	waitCh := make(chan struct{}, 1)
	idCh := make(chan string, 1)
	errCh := make(chan string, 1)
	defer close(lineCh)
	defer close(waitCh)
	defer close(idCh)
	defer close(errCh)

	go cmd.waitID(idCh, errCh, lineCh, waitCh)
	lineCh <- "testline1"
	lineCh <- "[network1] STARTED. INSTANCE [12345]"
	lineCh <- "testline2"

	timer := time.NewTimer(time.Second * 3)
	select {
	case id := <-idCh:
		if id != "12345" {
			t.Errorf("id => %s, want %s", id, "12345")
		}
	case errMsg := <-errCh:
		t.Errorf("Unexpected err received: %s", errMsg)
	case <-timer.C:
		t.Errorf("Test timeout.")
	}
}

func TestWaitID_ProcessEnd(t *testing.T) {
	cmd := new(Command)
	lineCh := make(chan string, 10)
	waitCh := make(chan struct{}, 1)
	idCh := make(chan string, 1)
	errCh := make(chan string, 1)
	defer close(lineCh)
	defer close(waitCh)
	defer close(idCh)
	defer close(errCh)

	go cmd.waitID(idCh, errCh, lineCh, waitCh)
	lineCh <- "testline1"
	lineCh <- "testline2"
	lineCh <- "testline3"
	time.Sleep(time.Millisecond * 10)
	waitCh <- struct{}{}

	timer := time.NewTimer(time.Second * 3)
	select {
	case id := <-idCh:
		t.Errorf("Unexpected id received: %s", id)
	case errMsg := <-errCh:
		if errMsg != "testline1\ntestline2\ntestline3\n" {
			t.Errorf("errMsg => %s, want %s", errMsg, "testline1\ntestline2\ntestline3\n")
		}
	case <-timer.C:
		t.Errorf("Test timeout.")
	}
}
