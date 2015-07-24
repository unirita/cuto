package network

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"cuto/util"
)

// Command represents a master command which executes realtime network.
type Command struct {
	cmd         *exec.Cmd
	networkName string
}

// NewCommand creates Command object with unique network name.
// Network name is generate from realtimeName and current timestamp.
func NewCommand(realtimeName string) *Command {
	c := new(Command)
	timestamp := time.Now().Format("20060102150405")
	if realtimeName == "" {
		c.networkName = fmt.Sprintf("realtime_%s", timestamp)
	} else {
		c.networkName = fmt.Sprintf("realtime_%s_%s", realtimeName, timestamp)
	}

	masterPath := filepath.Join(util.GetRootPath(), "bin", "master")
	configPath := filepath.Join(util.GetRootPath(), "bin", "master.ini")
	c.cmd = exec.Command(masterPath, "-n", c.networkName, "-s", "-c", configPath)
	return c
}

// GetNetworkName returns network name.
func (c *Command) GetNetworkName() string {
	return c.networkName
}

// Run runs the master command and gets its instance id.
func (c *Command) Run() (string, error) {
	stdoutReader, err := c.cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := c.cmd.Start(); err != nil {
		return "", err
	}

	lineCh := make(chan string, 1)
	waitCh := make(chan struct{}, 1)
	idCh := make(chan string, 1)
	errCh := make(chan string, 1)

	go c.monitorStdout(lineCh, stdoutReader)
	go c.waitProcess(waitCh)
	go c.waitID(idCh, errCh, lineCh, waitCh)

	select {
	case id := <-idCh:
		return id, nil
	case errMsg := <-errCh:
		return "", fmt.Errorf("Master error: %s", errMsg)
	}
}

// Release releases any resources associated with the master command process.
// It is recommended that call this function if you do not wait end of process.
func (c *Command) Release() error {
	return c.cmd.Process.Release()
}

func (c *Command) monitorStdout(lineCh chan<- string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lineCh <- line
		}
	}
}

func (c *Command) waitProcess(waitCh chan<- struct{}) {
	c.cmd.Wait()
	close(waitCh)
}

func (c *Command) waitID(idCh, errCh chan<- string, lineCh <-chan string, waitCh <-chan struct{}) {
	matcher := regexp.MustCompile(`INSTANCE \[\d+`)
	var lastLine string
	for {
		select {
		case lastLine = <-lineCh:
			id := matcher.FindString(lastLine)
			if id != "" {
				id = strings.Replace(id, "INSTANCE [", "", 1)
				idCh <- id
				return
			}
		case <-waitCh:
			errCh <- lastLine
			return
		}
	}
}
