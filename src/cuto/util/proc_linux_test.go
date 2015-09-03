package util

import (
	"os/exec"
	"strconv"
)

func createSleepCommand(second int) {
	return exec.Command("/bin/sleep", strconv.Itoa(second))
}
