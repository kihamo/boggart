package v1

import (
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

const (
	displayModeTariff1 uint8 = 1 << iota
	displayModeTariff2
	displayModeTariff3
	displayModeTariff4
	displayModeAmount
	displayModePower
	displayModeTime
	displayModeDate
)

const (
	displayModeTariffSchedule uint8 = 1 << iota
	displayModeUIF
	displayModeReactiveEnergy
	displayModeMaximumResets
	displayModeWorkingTime
	displayModeBatteryLifetime
	displayModePowerLimit
	displayModeEnergyLimit

	MaxEventsIndex = 0x3F
	CurrentMonth   = 0x0F

	MaximumPower    = 0x0
	MaximumAmperage = 0x1
	MaximumVoltage  = 0x2
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

func (m *MercuryV1) Invoke(request *Request) (*Response, error) {
	if request.address == 0 {
		request = request.WithAddress(m.options.address)
	}

	if request.address == 0 {
		return nil, errors.New("device address is empty")
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
	if response.address != request.address {
		return nil, fmt.Errorf("error ADDR of response packet %X have %X want %X",
			data, response.address, request.address)
	}

	return response, nil
}
