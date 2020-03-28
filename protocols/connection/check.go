// +build !windows,!appengine

package connection

import (
	"errors"
	"io"
	"net"
	"syscall"
)

var errUnexpectedRead = errors.New("unexpected read from socket")

func Check(c net.Conn) error {
	var (
		n    int
		err  error
		buff [1]byte
	)

	sconn, ok := c.(syscall.Conn)

	if !ok {
		return nil
	}

	rc, err := sconn.SyscallConn()

	if err != nil {
		return err
	}

	rerr := rc.Read(func(fd uintptr) bool {
		n, err = syscall.Read(int(fd), buff[:])
		return true
	})

	switch {
	case rerr != nil:
		return rerr
	case n == 0 && err == nil:
		return io.EOF
	case n > 0:
		return errUnexpectedRead
	case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
		return nil
	default:
		return err
	}
}
