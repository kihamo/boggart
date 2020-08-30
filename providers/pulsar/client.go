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
	conn.ApplyOptions(connection.WithReadCheck(ReadCheck))

	client := &HeatMeter{
		connection: conn,
		options:    defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&client.options)
	}

	return client
}

func (d *HeatMeter) Invoke(request *Packet) (*Packet, error) {
	if len(request.address) == 0 {
		request = request.WithAddress(d.options.address)
	}

	requestData, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	responseData, err := d.connection.Invoke(requestData)
	if err != nil {
		return nil, err
	}

	response := NewPacket()
	if err = response.UnmarshalBinary(responseData); err != nil {
		return nil, err
	}

	// check ADDR
	if !bytes.Equal(response.Address(), request.Address()) {
		return nil, errors.New(
			"error ADDR of response packet " +
				hex.EncodeToString(responseData) + " have " +
				hex.EncodeToString(response.Address()) + " want " +
				hex.EncodeToString(request.Address()))
	}

	// check ID
	if !bytes.Equal(response.ID(), request.ID()) {
		return nil, errors.New(
			"error ID of response packet " +
				hex.EncodeToString(responseData) + " have " +
				hex.EncodeToString(response.ID()) + " want " +
				hex.EncodeToString(request.ID()))
	}

	// check error
	if response.Function() == FunctionBadCommand {
		return nil, fmt.Errorf("returns error code %v", response.ErrorCode())
	}

	return response, nil
}
