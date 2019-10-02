package connection

import (
	"bytes"
	"io"
	"sync"
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
	if locker, ok := i.Conn.(sync.Locker); ok {
		locker.Lock()
		defer locker.Unlock()
	}

	_, err := i.Conn.Write(request)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	bufferTmp := make([]byte, bufferSize)

	for {
		n, err := i.Conn.Read(bufferTmp)

		if n > 0 {
			buffer.Write(bufferTmp[:n])
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
