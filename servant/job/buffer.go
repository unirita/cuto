package job

import (
	"bytes"
	"io"
	"os"
)

type StdoutBuffer struct {
	enableOutput bool
	bytes.Buffer
}

func NewStdoutBuffer(enableOutput bool) *StdoutBuffer {
	return &StdoutBuffer{enableOutput: enableOutput}
}

func (b *StdoutBuffer) Write(p []byte) (int, error) {
	if b.enableOutput {
		os.Stdout.Write(p)
	}
	return b.Buffer.Write(p)
}

type OutputPipeBuffer struct {
	enableOutput bool
	bytes.Buffer
}

func NewOutputPipeBuffer(enableOutput bool) *OutputPipeBuffer {
	return &OutputPipeBuffer{enableOutput: enableOutput}
}

func (b *OutputPipeBuffer) ReadPipe(pStdout io.ReadCloser, pStderr io.ReadCloser) error {
	reader := io.MultiReader(pStdout, pStderr)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return err
		}
		if b.enableOutput {
			os.Stdout.Write(buf[:n])
		}
		b.Buffer.Write(buf[:n])
	}
	return nil
}
