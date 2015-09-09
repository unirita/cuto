package util

import (
	"os/exec"
	"strconv"
)

func createSleepCommand(second int) *exec.Cmd {
	second++
	return exec.Command("ping", "-n", strconv.Itoa(second), "localhost")
}
