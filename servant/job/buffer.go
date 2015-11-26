package job

import (
	"bytes"
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
