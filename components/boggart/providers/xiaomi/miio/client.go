package miio

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal/packet"
)

type Client struct {
	io.Closer

	conn           *internal.Connection
	connOnce       sync.Once
	packetsCounter uint32

	address  string
	deviceId []byte
	token    string

	stampDiff time.Duration
}

func NewClient(address, token string) *Client {
	return &Client{
		packetsCounter: uint32(time.Now().Unix()),
		address:        address,
		token:          token,
	}
}

func (p *Client) lazyConnect() (conn *internal.Connection, err error) {
	p.connOnce.Do(func() {
		conn, err = internal.NewConnection(p.address)
		if err == nil {
			p.conn = conn
		}
	})

	return p.conn, nil
}

func (p *Client) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}

func (p *Client) Hello() (packet.Packet, error) {
	conn, err := p.lazyConnect()
	if err != nil {
		return nil, err
	}

	request := packet.NewHello()
	response := packet.NewBase()

	err = conn.Invoke(request, response)
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

	conn, err := p.lazyConnect()
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

	fmt.Println(string(body))

	request, err := packet.NewCrypto(deviceID, p.token)
	if err != nil {
		return err
	}
	request.SetBody(body)
	request.SetStamp(time.Now().Add(p.stampDiff))

	err = conn.Invoke(request, response)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	fmt.Println(string(response.Body()))

	err = json.Unmarshal(response.Body(), &result)
	return err
}
