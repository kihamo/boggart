package connection

import (
	"io"
	"net"
	"time"
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

func SetDeadline(duration time.Duration, f func(t time.Time) error) error {
	var deadline time.Time

	if duration > 0 {
		deadline = time.Now().Add(duration)
	} else {
		deadline = time.Time{}
	}

	return f(deadline)
}
