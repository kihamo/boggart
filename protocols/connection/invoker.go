package connection

import (
	"bytes"
	"io"
)

const (
	bufferSize = 512
)

type Invoker interface {
	Conn

	Invoke(request []byte) (response []byte, err error)
}

type invoker struct {
	Conn
}

func NewInvoker(conn Conn) Invoker {
	if i, ok := conn.(Invoker); ok {
		return i
	}

	return &invoker{
		Conn: conn,
	}
}

func (i *invoker) Invoke(request []byte) ([]byte, error) {
	_, err := i.Write(request)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	for {
		b := make([]byte, bufferSize)
		n, err := i.Read(b)

		if n > 0 {
			buffer.Write(b[:n])
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return buffer.Bytes(), nil
}
