package v3

import (
	"fmt"

	"github.com/kihamo/boggart/protocols/connection"
)

type MercuryV3 struct {
	invoker connection.Invoker
	options options
}

func New(conn connection.Conn, opts ...Option) *MercuryV3 {
	m := &MercuryV3{
		invoker: connection.NewInvoker(conn),
		options: defaultOptions(),
	}

	for _, opt := range opts {
		opt.apply(&m.options)
	}

	return m
}

func (m *MercuryV3) Invoke(request *Request) (response *Response, err error) {
	err = m.ChannelOpen(m.options.accessLevel, m.options.password)
	if err != nil {
		return response, err
	}

	response, err = m.InvokeRaw(request)
	if err == nil {
		err = m.ChannelClose()
	}

	return response, err
}

func (m *MercuryV3) InvokeRaw(request *Request) (*Response, error) {
	if request.address == 0 {
		request = request.WithAddress(m.options.address)
	}

	//if request.address == 0 {
	//	return nil, errors.New("device address is empty")
	//}

	data, err := m.invoker.Invoke(request.Bytes())
	if err != nil {
		return nil, err
	}

	return ParseResponse(data)
}

func (m *MercuryV3) Raw() error {
	bwri := uint8(0x8<<4) | uint8(PhaseNumber1)
	resp, err := m.AuxiliaryParameters(bwri)

	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}
