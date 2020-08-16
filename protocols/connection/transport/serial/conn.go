package serial

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/goburrow/serial"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/protocols/connection/transport"
)

type Serial struct {
	options   options
	once      *atomic.Once
	port      serial.Port
	portMutex sync.Mutex // внутри port поле fd не защищено от race
}

func New(options ...Option) *Serial {
	c := &Serial{
		options: DefaultOptions(),
	}

	for _, option := range options {
		option.apply(&c.options)
	}

	return c
}

func (s *Serial) Dial() (_ transport.Transport, err error) {
	w := &Serial{
		options: s.options,
	}

	w.port, err = serial.Open(&s.options.Config)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *Serial) Read(p []byte) (n int, err error) {
	if s.port == nil {
		return -1, errors.New("connection isn't init")
	}

	buffer := bytes.NewBuffer(nil)

	for {
		s.portMutex.Lock()
		n, err = s.port.Read(p)
		s.portMutex.Unlock()

		if n > 0 {
			buffer.Write(p[:n])
		}

		if err != nil {
			if err == serial.ErrTimeout {
				err = io.EOF
			}

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

	return n, err
}

func (s *Serial) Write(p []byte) (n int, err error) {
	if s.port == nil {
		return -1, errors.New("connection isn't init")
	}

	s.portMutex.Lock()
	defer s.portMutex.Unlock()

	return s.port.Write(p)
}

func (s *Serial) Close() error {
	if s.port != nil {
		s.portMutex.Lock()
		defer s.portMutex.Unlock()

		return s.port.Close()
	}

	return nil
}

func (s *Serial) Options() map[string]interface{} {
	return s.options.Map()
}
