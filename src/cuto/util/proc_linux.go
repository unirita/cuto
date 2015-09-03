package util

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func isProcessExists(pid int) bool {
	ptn := fmt.Sprintf(`^\s*%d$`, pid)

	pgrepCmd := exec.Command("ps", "-eo", "pid")
	grepCmd := exec.Command("grep", ptn)
	wcCmd := exec.Command("wc", "-l")

	r1, w1 := io.Pipe()
	pgrepCmd.Stdout = w1
	grepCmd.Stdin = r1

	r2, w2 := io.Pipe()
	grepCmd.Stdout = w2
	wcCmd.Stdin = r2

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
	w1.Close()
	grepCmd.Wait()
	w2.Close()
	wcCmd.Wait()

	s := strings.Trim(b.String(), " \t\r\n")
	if s == "0" {
		return false
	}
	return true
}
