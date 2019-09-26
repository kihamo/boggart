package pulsar

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

type HeatMeter struct {
	invoker connection.Invoker
	options options
}

func New(conn connection.Conn, opts ...Option) *HeatMeter {
	client := &HeatMeter{
		invoker: connection.NewInvoker(conn),
		options: defaultOptions(),
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

	data, err := d.invoker.Invoke(request.Bytes())
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

	// check ID
	if !bytes.Equal(response.Id, request.Id) {
		return nil, errors.New(
			"error ID of response packet have " +
				hex.EncodeToString(response.Id) + " want " +
				hex.EncodeToString(request.Id))
	}

	// check error
	if response.Function == FunctionBadCommand {
		return nil, fmt.Errorf("returns error code #%d", response.ErrorCode)
	}

	return response, nil
}
