package nis

// http://www.apcupsd.org/manual/manual.html#nis-network-server-protocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
)

type reader struct {
	address string
	command string
}

func (r *reader) Reader(ctx context.Context) (io.Reader, error) {
	var d net.Dialer

	connect, err := d.DialContext(ctx, "tcp", r.address)
	if err != nil {
		return nil, err
	}
	defer connect.Close()

	lengthBuffer := make([]byte, 2)

	// request
	binary.BigEndian.PutUint16(lengthBuffer, uint16(len(r.command)))

	if _, err = connect.Write(append(lengthBuffer, []byte(r.command)...)); err != nil {
		return nil, err
	}

	// response
	response := bytes.NewBuffer(nil)

	for {
		if _, err := io.ReadFull(connect, lengthBuffer); err != nil {
			return nil, err
		}

		chunkLength := binary.BigEndian.Uint16(lengthBuffer)
		if chunkLength == 0 {
			break
		}

		chunk := make([]byte, chunkLength)
		if _, err = io.ReadFull(connect, chunk[:chunkLength]); err != nil {
			return nil, err
		}

		response.Write(chunk)
	}

	return response, nil
}
