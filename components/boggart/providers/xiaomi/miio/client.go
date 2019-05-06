package miio

import (
	"encoding/json"
	"io"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal/packet"
)

type Client struct {
	io.Closer

	conn           *internal.Connection
	packetsCounter uint32

	deviceId []byte
	token    string

	stampDiff time.Duration
}

func NewClient(address, token string) (*Client, error) {
	conn, err := internal.NewConnection(address)
	if err != nil {
		return nil, err
	}

	return &Client{
		packetsCounter: 6,
		token:          token,
		conn:           conn,
	}, nil
}

func (p *Client) Close() error {
	return p.conn.Close()
}

func (p *Client) Hello() (packet.Packet, error) {
	request := packet.NewHello()
	response := packet.NewBase()

	err := p.conn.Invoke(request, response)
	if err != nil {
		return nil, err
	}

	p.stampDiff = time.Now().Sub(response.Stamp())
	return response, nil
}

func (p *Client) DeviceID() ([]byte, error) {
	if p.deviceId == nil {
		response, err := p.Hello()
		if err != nil {
			return nil, err
		}

		p.deviceId = response.DeviceID()
	}

	return p.deviceId, nil
}

func (p *Client) Send(method string, params interface{}, result interface{}) error {
	if params == nil {
		params = []interface{}{}
	}

	body, err := json.Marshal(Request{
		ID:     atomic.AddUint32(&p.packetsCounter, 1),
		Method: method,
		Params: params,
	})

	if err != nil {
		return err
	}

	deviceID, err := p.DeviceID()
	if err != nil {
		return err
	}

	var response packet.Packet

	if result != nil {
		response, err = packet.NewCrypto(deviceID, p.token)
		if err != nil {
			return err
		}
	}

	request, err := packet.NewCrypto(deviceID, p.token)
	if err != nil {
		return err
	}
	request.SetBody(body)
	request.SetStamp(time.Now().Add(p.stampDiff))

	err = p.conn.Invoke(request, response)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	err = json.Unmarshal(response.Body(), &result)
	return err
}
