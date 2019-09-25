package v1

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/providers/mercury"
)

const (
	displayModeTariff1 = 1 << iota
	displayModeTariff2
	displayModeTariff3
	displayModeTariff4
	displayModeAmount
	displayModePower
	displayModeTime
	displayModeDate
)

type MercuryV1 struct {
	connection mercury.Connection
	options    options
}

func New(connection mercury.Connection, opts ...Option) *MercuryV1 {
	m := &MercuryV1{
		connection: connection,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	return m
}

func (m *MercuryV1) Request(request *Request) (*Response, error) {
	if len(request.Address) == 0 {
		request.Address = m.options.address
	}

	data, err := m.connection.Invoke(request.Bytes())
	if err != nil {
		return nil, err
	}

	response, err := ParseResponse(data)
	if err != nil {
		return nil, err
	}

	// check ADDR
	if !bytes.Equal(response.Address, request.Address) {
		return nil, errors.New(
			"error ADDR of response packet have " +
				hex.EncodeToString(response.Address) + " want " +
				hex.EncodeToString(request.Address))
	}

	return response, nil
}
