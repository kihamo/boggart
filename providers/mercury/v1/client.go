package v1

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/protocols/connection"
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
	invoker connection.Invoker
	options options
}

func New(conn connection.Conn, opts ...Option) *MercuryV1 {
	m := &MercuryV1{
		invoker: connection.NewInvoker(conn),
		options: defaultOptions(),
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

	data, err := m.invoker.Invoke(request.Bytes())
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
			"error ADDR of response packet " +
				hex.EncodeToString(data) + " have " +
				hex.EncodeToString(response.Address) + " want " +
				hex.EncodeToString(request.Address))
	}

	return response, nil
}
