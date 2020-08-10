package transport

import (
	"io"
)

type Transport interface {
	io.ReadWriteCloser
	Dial() (Transport, error)
	Options() map[string]interface{}
}
