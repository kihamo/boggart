package miio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	a "sync/atomic"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/providers/xiaomi/miio/internal"
	"github.com/kihamo/boggart/providers/xiaomi/miio/internal/packet"
)

const (
	DefaultTimeout = time.Second * 5
)

type Client struct {
	io.Closer

	packetsCounter uint32
	dump           uint32
	stampDiff      int64

	conn     *internal.Connection
	connOnce *atomic.Once

	address  string
	deviceID []byte
	token    string
}

func NewClient(address, token string) *Client {
	return &Client{
		address:  address,
		token:    token,
		connOnce: new(atomic.Once),
	}
}

func (p *Client) lazyConnect(ctx context.Context) (_ *internal.Connection, err error) {
	p.connOnce.Do(func() {
		p.conn, err = internal.NewConnection(p.address)
		if err == nil {
			// p.SetDump(true)

			// начинаем сессию hello пакетом
			_, err = p.Hello(ctx)
		}
	})

	if err != nil {
		p.connOnce.Reset()
	}

	return p.conn, err
}

func (p *Client) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}

func (p *Client) LocalAddr() (*net.UDPAddr, error) {
	connect, err := p.lazyConnect(context.Background())
	if err != nil {
		return nil, err
	}

	return connect.LocalAddr(), nil
}

func (p *Client) SetDump(enabled bool) {
	if enabled {
		a.StoreUint32(&p.dump, 1)
	} else {
		a.StoreUint32(&p.dump, 0)
	}
}

func (p *Client) SetPacketsCounter(count uint32) {
	a.StoreUint32(&p.packetsCounter, count)
}

func (p *Client) Hello(ctx context.Context) (packet.Packet, error) {
	conn, err := p.lazyConnect(ctx)
	if err != nil {
		return nil, err
	}

	request := packet.NewHello()
	response := packet.NewBase()

	err = conn.Invoke(ctx, request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *Client) DeviceID(ctx context.Context) ([]byte, error) {
	if p.deviceID == nil {
		response, err := p.Hello(ctx)
		if err != nil {
			return nil, err
		}

		p.deviceID = response.DeviceID()
	}

	return p.deviceID, nil
}

func (p *Client) Send(ctx context.Context, method string, params interface{}, result interface{}) error {
	if params == nil {
		params = []interface{}{}
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)

		defer cancel()
	}

	requestID := a.AddUint32(&p.packetsCounter, 1)
	body, err := json.Marshal(Request{
		ID:     requestID,
		Method: method,
		Params: params,
	})

	if err != nil {
		return err
	}

	deviceID, err := p.DeviceID(ctx)
	if err != nil {
		return err
	}

	conn, err := p.lazyConnect(ctx)
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

	var diff time.Duration

	diff = time.Duration(a.LoadInt64(&p.stampDiff))
	request.SetStamp(time.Now().Add(diff))

	dump := a.LoadUint32(&p.dump) == 1
	if dump {
		log.Printf("Request #%d raw: "+request.Dump(), requestID)
		log.Printf("Request #%d body: "+string(body), requestID)
	}

	err = conn.Invoke(ctx, request, response)

	if err != nil {
		return err
	}

	if dump {
		log.Println("Response raw: " + response.Dump())
		log.Println("Response body: " + string(response.Body()))
	}

	if result == nil {
		return nil
	}

	diff = time.Since(response.Stamp())
	a.StoreInt64(&p.stampDiff, int64(diff))

	var responseError ResponseError
	if err = json.Unmarshal(response.Body(), &responseError); err == nil && len(responseError.Error.Message) > 0 {
		return errors.New(responseError.Error.Message)
	}

	var responseUnknown ResponseUnknownMethod
	if err = json.Unmarshal(response.Body(), &responseUnknown); err == nil && responseUnknown.Result == "unknown_method" {
		return errors.New("unknown method")
	}

	var responseDefault Response

	// nolint:wsl
	if err = json.Unmarshal(response.Body(), &responseDefault); err == nil && responseDefault.ID != requestID {
		// TODO: если requestID > responseDefault.ID то можно еще делать попытки вычитывания, обычно такое бывает
		// из-за того что контекст прерывается быстрее чем девай отвечает и следующий реквест получает предыдущий респонс

		return errors.New("response ID could not be verified. Expected " +
			strconv.FormatUint(uint64(requestID), 10) +
			", got " +
			strconv.FormatUint(uint64(responseDefault.ID), 10))
	}

	err = json.Unmarshal(response.Body(), &result)

	return err
}
