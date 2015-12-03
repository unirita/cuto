// +build darwin linux

package util

import (
	"os/exec"
	"strconv"
)

func createSleepCommand(second int) *exec.Cmd {
	return exec.Command("/bin/sleep", strconv.Itoa(second))
}
