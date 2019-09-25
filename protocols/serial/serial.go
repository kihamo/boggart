package serial

import (
	"bytes"
	"fmt"
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
		multiRequestsConnections[conn.options.target] = &sync.Mutex{}
		multiRequestsMutex.Unlock()
	}

	return conn
}

func (c *Serial) lockWrapper(f func(s.Port) error) error {
	multiRequestsMutex.Lock()
	lock, ok := multiRequestsConnections[c.options.target]
	multiRequestsMutex.Unlock()

	if ok {
		lock.Lock()
		defer lock.Unlock()
	}

	port, err := s.Open(&s.Config{
		BaudRate: c.options.baudRate,
		DataBits: c.options.dataBits,
		StopBits: c.options.stopBits,
		Parity:   c.options.parity,
		Address:  c.options.target,
		Timeout:  c.options.timeout,
	})

	if err != nil {
		return err
	}
	defer port.Close()

	return f(port)
}

func (c *Serial) Read(p []byte) (n int, err error) {
	err = c.lockWrapper(func(port s.Port) (e error) {
		buffer := bytes.NewBuffer(nil)

		for {
			b := make([]byte, len(p))
			rn, re := port.Read(b)

			if rn > 0 {
				buffer.Write(b[:rn])
			}

			if re != nil {
				if re.Error() == "serial: timeout" {
					e = io.EOF
					break
				}

				e = re
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

		return e
	})

	return n, err
}

func (c *Serial) Write(p []byte) (n int, err error) {
	err = c.lockWrapper(func(port s.Port) (e error) {
		n, e = port.Write(p)
		return e
	})

	return n, err
}

func (c *Serial) ReadWrite(reader io.Reader, writer io.Writer) error {
	return c.lockWrapper(func(port s.Port) (e error) {
		readSerialDone := make(chan struct{}, 1)
		readSerial := make(chan struct{}, 1)
		errors := make(chan error, 1)

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
		if err != nil && err.Error() == "serial: timeout" {
			return nil
		}

		return err
	})
}

func (c *Serial) Invoke(request []byte) (response []byte, err error) {
	fmt.Println("Invoke")

	err = c.lockWrapper(func(port s.Port) (e error) {
		if _, e = port.Write(request); e != nil {
			return e
		}

		buffer := bytes.NewBuffer(nil)

		for {
			b := make([]byte, bufferSize)
			n, e := port.Read(b)
			if e != nil {
				break
			}

			if n != 0 {
				buffer.Write(b[:n])
			}
		}

		response = buffer.Bytes()
		return e
	})

	return response, err
}
