package util

import (
	"bytes"
	"io"
	"os/exec"
)

func IsProcessExists(pid int) bool {
	pgrepCmd := exec.Command("pgrep", "master")
	grepCmd := exec.Command("grep", "^"+pid+"$")
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

	if b.String() == "0" {
		return true
	}
	return false
}
