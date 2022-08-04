package modbus

import (
	"bytes"
	"io"
)

type logger struct {
	io.Writer

	f func(string)
}

func NewLogger(f func(string)) io.Writer {
	return &logger{
		f: f,
	}
}

func (l logger) Write(p []byte) (n int, err error) {
	l.f(string(bytes.TrimSpace(p)))

	return len(p), nil
}
