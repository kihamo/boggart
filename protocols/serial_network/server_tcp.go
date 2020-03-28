package serialnetwork

import (
	"net"
	"sync"
)

type TCPServer struct {
	network string
	address string
	serial  ReadWriter

	lock     sync.Mutex
	listener net.Listener
}

func NewTCPServer(network, address string, serial ReadWriter) *TCPServer {
	return &TCPServer{
		network: network,
		address: address,
		serial:  serial,
	}
}

func (s *TCPServer) ListenAndServe() error {
	listen, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.listener = listen
	s.lock.Unlock()

	for {
		conn, err := listen.Accept()
		if err != nil {
			if x, ok := err.(*net.OpError); ok && x.Op == "accept" {
				break
			}

			// TODO: log
			continue
		}

		go func() {
			defer conn.Close()
			_ = s.serial.ReadWrite(conn, conn)
		}()
	}

	return nil
}

func (s *TCPServer) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
