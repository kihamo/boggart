package connection

import (
	"io"
	"net"
)

type Conn interface {
	io.ReadWriteCloser
}

func IsTimeout(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}
