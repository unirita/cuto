package util

import (
	"testing"
	"time"
)

func TestIsProcessExists_ReturnTrueWhenExists(t *testing.T) {
	cmd := createSleepCommand(5)
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Cannot start sleep process.")
	}
	defer cmd.Process.Kill()

	if !IsProcessExists(cmd.Process.Pid) {
		t.Error("IsProcessExists() returns false.")
	}
}

func TestIsProcessExists_ReturnFalseWhenNotExists(t *testing.T) {
	cmd := createSleepCommand(0)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Cannot start sleep process.")
	}

	time.Sleep(time.Second)

	if IsProcessExists(cmd.Process.Pid) {
		t.Error("IsProcessExists() returns true.")
	}
}
