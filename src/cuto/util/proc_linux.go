package util

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func IsProcessExists(pid int) bool {
	ptn := fmt.Sprintf(`^\s*%d$`, pid)

	pgrepCmd := exec.Command("ps", "-eo", "pid")
	grepCmd := exec.Command("grep", ptn)
	wcCmd := exec.Command("wc", "-l")

	r1, w1 := io.Pipe()
	pgrepCmd.Stdout = w1
	grepCmd.Stdin = r1
	defer w1.Close()
	defer r1.Close()

	r2, w2 := io.Pipe()
	grepCmd.Stdout = w2
	wcCmd.Stdin = r2
	defer w2.Close()
	defer r2.Close()

	b := new(bytes.Buffer)
	wcCmd.Stdout = b

	if err := pgrepCmd.Start(); err != nil {
		panic(err)
	}
	if err := grepCmd.Start(); err != nil {
		panic(err)
	}
	if err := wcCmd.Start(); err != nil {
		panic(err)
	}
	pgrepCmd.Wait()
	grepCmd.Wait()
	wcCmd.Wait()

	s := strings.Trim(b.String(), " \t\r\n")
	if s == "0" {
		return false
	}
	return true
}
