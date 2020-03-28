package serial_network

import (
	"bytes"
	"net"
	"sync"
)

const (
	maxBufferSize = 1024
)

type UDPServer struct {
	network string
	address string
	serial  ReadWriter

	lock     sync.Mutex
	listener net.PacketConn
}

func NewUDPServer(network, address string, serial ReadWriter) *UDPServer {
	return &UDPServer{
		network: network,
		address: address,
		serial:  serial,
	}
}

func (s *UDPServer) ListenAndServe() error {
	listen, err := net.ListenPacket(s.network, s.address)
	if err != nil {
		return err
	}

	s.lock.Lock()
	s.listener = listen
	s.lock.Unlock()

	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := listen.ReadFrom(buffer)
		if err != nil {
			continue
		}

		if n < 1 {
			continue
		}

		// TODO: deadline configs

		go func() {
			reader := bytes.NewReader(buffer[:n])
			writer := bytes.NewBuffer(nil)

			err := s.serial.ReadWrite(reader, writer)
			if err == nil {
				if writer.Len() > 0 {
					_, _ = listen.WriteTo(writer.Bytes(), addr)
				}
			}

			// TODO: log
		}()
	}
}

func (s *UDPServer) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
