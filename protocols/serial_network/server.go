package serial_network

import (
	"io"
)

type ReadWriter interface {
	ReadWrite(reader io.Reader, writer io.Writer) error
}

type Server interface {
	ListenAndServe() error
	Close() error
}
