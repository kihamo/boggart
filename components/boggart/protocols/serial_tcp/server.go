package serial_tcp

import (
	"fmt"
	"io"
	"net"
)

type ReadWriter interface {
	ReadWrite(reader io.Reader, writer io.Writer) error
}

type Server struct {
	listener net.Listener
	address  string

	serial ReadWriter
}

func NewServer(address string, serial ReadWriter) *Server {
	return &Server{
		address: address,
		serial:  serial,
	}
}

func (s *Server) ListenAndServe() error {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			// TODO: log
			fmt.Println("Error ", err.Error())
		}

		go func() {
			defer conn.Close()
			s.serial.ReadWrite(conn, conn)
		}()
	}
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
