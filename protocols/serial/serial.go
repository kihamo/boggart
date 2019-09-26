package serial

import (
	"bytes"
	"io"
	"sync"

	s "github.com/goburrow/serial"
)

const (
	bufferSize = 512
)

var (
	multiRequestsMutex       sync.RWMutex
	multiRequestsConnections = make(map[string]*sync.Mutex)
)

type Serial struct {
	options options
}

func Dial(opts ...Option) *Serial {
	conn := &Serial{
		options: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&conn.options)
	}

	if !conn.options.allowMultiRequest {
		multiRequestsMutex.Lock()
		multiRequestsConnections[conn.options.Config.Address] = &sync.Mutex{}
		multiRequestsMutex.Unlock()
	}

	return conn
}

func (c *Serial) Lock() {
	multiRequestsMutex.Lock()
	lock, ok := multiRequestsConnections[c.options.Config.Address]
	multiRequestsMutex.Unlock()

	if ok {
		lock.Lock()
	}
}

func (c *Serial) Unlock() {
	multiRequestsMutex.Lock()
	lock, ok := multiRequestsConnections[c.options.Config.Address]
	multiRequestsMutex.Unlock()

	if ok {
		lock.Unlock()
	}
}

func (c *Serial) Read(p []byte) (n int, err error) {
	port, err := s.Open(&c.options.Config)
	if err == nil {
		buffer := bytes.NewBuffer(nil)

		for {
			b := make([]byte, len(p))
			rn, re := port.Read(b)

			if rn > 0 {
				buffer.Write(b[:rn])
			}

			if re != nil {
				if re == s.ErrTimeout {
					err = io.EOF
					break
				}

				err = re
				break
			}

			if len(p) <= buffer.Len() {
				break
			}
		}

		if buffer.Len() > len(p) {
			n = len(p)
		} else {
			n = buffer.Len()
		}

		copy(p, buffer.Bytes()[:n])
	}

	return n, err
}

func (c *Serial) Write(p []byte) (n int, err error) {
	port, err := s.Open(&c.options.Config)
	if err == nil {
		n, err = port.Write(p)
	}

	return n, err
}

func (c *Serial) ReadWrite(reader io.Reader, writer io.Writer) error {
	port, e := s.Open(&c.options.Config)
	if e != nil {
		return e
	}

	readSerialDone := make(chan struct{}, 1)
	readSerial := make(chan struct{}, 1)
	errors := make(chan error, 1)

	c.Lock()
	defer c.Unlock()

	// from SERIAL to WRITER
	go func() {
		for {
			select {
			case <-readSerial:
				for {
					serialBuffer := make([]byte, bufferSize)

					n, err := port.Read(serialBuffer)
					if err != nil {
						// любая ошибка при чтении из порта делает не возможную работу с ним в рамках этой сессии
						errors <- err
						return
					}

					if n > 0 {
						_, err = writer.Write(serialBuffer[:n])
						if err != nil {
							return
						}
					}
				}

			case <-readSerialDone:
				return
			}
		}
	}()

	// from READER to SERIAL
	go func() {
		for {
			readerBuffer := make([]byte, bufferSize)

			n, err := reader.Read(readerBuffer)
			if err != nil {
				return
			}

			if n < 1 {
				continue
			}

			_, err = port.Write(readerBuffer[:n])
			if err != nil {
				errors <- err
				return
			}

			readSerial <- struct{}{}
		}
	}()

	defer close(readSerialDone)

	err := <-errors
	if err != nil && err == s.ErrTimeout {
		return nil
	}

	return err
}
