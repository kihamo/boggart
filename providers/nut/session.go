package nut

import (
	"bytes"
	"errors"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
	protocol "github.com/kihamo/boggart/protocols/connection"
)

var (
	prefixError = []byte("ERR ")
	prefixBegin = []byte("BEGIN ")
	prefixEnd   = []byte("END ")
	separator   = []byte("\n")
	resultOK    = []byte("OK")

	readBufferPool sync.Pool
)

func init() {
	readBufferPool.New = func() interface{} {
		buf := make([]byte, 512)
		return &buf
	}
}

type Session struct {
	starter    int32
	connection protocol.Invoker

	dsn      string
	username string
	password string
}

func NewSession(dsn, username, password string) *Session {
	session := &Session{
		dsn:      dsn,
		username: username,
		password: password,
	}

	runtime.SetFinalizer(session, func(s *Session) {
		s.Close()
	})

	return session
}

func (s *Session) Start() error {
	if atomic.LoadInt32(&s.starter) == 1 {
		return nil
	}

	conn, err := protocol.New(s.dsn)

	if err != nil {
		return err
	}

	s.connection = conn.(protocol.Invoker)
	atomic.StoreInt32(&s.starter, 1)

	if s.username != "" {
		if err = s.Username(s.username); err != nil {
			return err
		}
	}

	if s.password != "" {
		if err = s.Password(s.password); err != nil {
			return err
		}
	}

	return err
}

func (s *Session) Close() (err error) {
	if atomic.LoadInt32(&s.starter) != 1 {
		return nil
	}

	err = s.Logout()

	if e := s.connection.Close(); e != nil {
		err = multierror.Append(err, e)
	}

	atomic.StoreInt32(&s.starter, 0)

	return err
}

func (s *Session) Request(cmd string) ([]byte, error) {
	if err := s.Start(); err != nil {
		return nil, err
	}

	response, err := s.connection.Invoke(append([]byte(cmd), separator...))
	if err != nil {
		return nil, err
	}

	if bytes.HasPrefix(response, prefixError) {
		return nil, ErrorByCode(string(bytes.TrimSpace(response[len(prefixError):])))
	}

	// multi line response
	if bytes.HasPrefix(response, prefixBegin) {
		multiCmd := bytes.Split(response[len(prefixBegin):], separator)[0]
		if string(multiCmd) != cmd {
			return nil, errors.New("response cmd could not be verified. Expected '" + cmd + "', got '" + string(multiCmd) + "'")
		}

		prefix := append(append(prefixBegin, []byte(cmd)...), separator...)
		suffix := append(append(prefixEnd, []byte(cmd)...), separator...)

		// cut prefix
		response = response[len(prefix):]

		buf := readBufferPool.Get().(*[]byte)
		defer readBufferPool.Put(buf)

		for {
			n, err := s.connection.Read(*buf)

			if n > 0 {
				response = append(response, (*buf)[:n]...)

				if bytes.HasSuffix(response, suffix) {
					// cut postfix
					response = response[:len(response)-len(suffix)-1]
					break
				}
			}

			if err != nil {
				if err == io.EOF {
					break
				}

				return nil, err
			}
		}
	}

	return bytes.TrimSuffix(response, separator), nil
}

func (s *Session) Call(cmd string) error {
	response, err := s.Request(cmd)
	if err != nil {
		return err
	}

	if !bytes.HasPrefix(response, resultOK) {
		return errors.New("bad response. Expected 'ОК', got '" + string(response) + "'")
	}

	return nil
}
