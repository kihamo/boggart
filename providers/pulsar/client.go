package pulsar

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

type HeatMeter struct {
	connection connection.Connection
	options    options
}

func New(conn connection.Connection, opts ...Option) *HeatMeter {
	conn.ApplyOptions(connection.WithGlobalLock(true))
	conn.ApplyOptions(connection.WithOnceInit(true))

	client := &HeatMeter{
		connection: conn,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&client.options)
	}

	return client
}

func (d *HeatMeter) Request(request *Request) (*Response, error) {
	if len(request.Address) == 0 {
		request.Address = d.options.address
	}

	data, err := d.connection.Invoke(request.Bytes())
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

	// check ID
	if !bytes.Equal(response.ID, request.ID) {
		return nil, errors.New(
			"error ID of response packet " +
				hex.EncodeToString(data) + " have " +
				hex.EncodeToString(response.ID) + " want " +
				hex.EncodeToString(request.ID))
	}

	// check error
	if response.Function == FunctionBadCommand {
		return nil, fmt.Errorf("returns error code #%d", response.ErrorCode)
	}

	return response, nil
}
