package xmeye

import (
	"io"
	"sync"
)

type stream struct {
	*connection

	claim   *Packet
	request func() error

	reader *io.PipeReader
	writer *io.PipeWriter
	once   sync.Once
}

func newStream(conn *connection, claim *Packet, request func() error) *stream {
	s := &stream{
		connection: conn,
		claim:      claim,
		request:    request,
	}

	s.reader, s.writer = io.Pipe()

	return s
}

func (s *stream) Close() (err error) {
	err = s.reader.Close()
	if err == nil {
		err = s.writer.Close()
	}

	return err
}

func (s *stream) Read(p []byte) (n int, err error) {
	s.once.Do(func() {
		err = s.send(s.claim)
		if err != nil {
			return
		}

		err = s.request()
		if err != nil {
			return
		}

		response, err := s.receive()
		if err != nil {
			return
		}

		if err := response.payload.Error(); err != nil {
			return
		}

		go func() {
			for {
				response, err := s.receive()
				if err != nil {
					s.writer.CloseWithError(err)
					return
				}

				response.payload.WriteTo(s.writer)

				if response.currentPacket == 0x01 {
					s.writer.Close()
					return
				}
			}
		}()
	})

	if err != nil {
		return -1, err
	}

	return s.reader.Read(p)
}
