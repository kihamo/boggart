package connection

import (
	"io"
)

type Conn interface {
	io.ReadWriteCloser
}
