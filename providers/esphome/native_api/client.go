package native_api

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	defaultPort   = 6053
	defaultClient = "GoLang Native API client"

	packetMagicByte byte = 0x00
)

var messageTypesByName = map[string]byte{
	"native_api.HelloRequest":                          0x01,
	"native_api.HelloResponse":                         0x02,
	"native_api.ConnectRequest":                        0x03,
	"native_api.ConnectResponse":                       0x04,
	"native_api.DisconnectRequest":                     0x05,
	"native_api.DisconnectResponse":                    0x06,
	"native_api.PingRequest":                           0x07,
	"native_api.PingResponse":                          0x08,
	"native_api.DeviceInfoRequest":                     0x09,
	"native_api.DeviceInfoResponse":                    0x0A,
	"native_api.ListEntitiesRequest":                   0x0B,
	"native_api.ListEntitiesBinarySensorResponse":      0x0C,
	"native_api.ListEntitiesCoverResponse":             0x0D,
	"native_api.ListEntitiesFanResponse":               0x0E,
	"native_api.ListEntitiesLightResponse":             0x0F,
	"native_api.ListEntitiesSensorResponse":            0x10,
	"native_api.ListEntitiesSwitchResponse":            0x11,
	"native_api.ListEntitiesTextSensorResponse":        0x12,
	"native_api.ListEntitiesDoneResponse":              0x13,
	"native_api.SubscribeStatesRequest":                0x14,
	"native_api.BinarySensorStateResponse":             0x15,
	"native_api.CoverStateResponse":                    0x16,
	"native_api.FanStateResponse":                      0x17,
	"native_api.LightStateResponse":                    0x18,
	"native_api.SensorStateResponse":                   0x19,
	"native_api.SwitchStateResponse":                   0x1A,
	"native_api.TextSensorStateResponse":               0x1B,
	"native_api.SubscribeLogsRequest":                  0x1C,
	"native_api.SubscribeLogsResponse":                 0x1D,
	"native_api.CoverCommandRequest":                   0x1E,
	"native_api.FanCommandRequest":                     0x1F,
	"native_api.LightCommandRequest":                   0x20,
	"native_api.SwitchCommandRequest":                  0x21,
	"native_api.SubscribeHomeassistantServicesRequest": 0x22,
	"native_api.HomeassistantServiceResponse":          0x23,
	"native_api.GetTimeRequest":                        0x24,
	"native_api.GetTimeResponse":                       0x25,
	"native_api.SubscribeHomeAssistantStatesRequest":   0x26,
	"native_api.SubscribeHomeAssistantStateResponse":   0x27,
	"native_api.HomeAssistantStateResponse":            0x28,
	"native_api.ListEntitiesServicesResponse":          0x29,
	"native_api.ExecuteServiceRequest":                 0x2A,
	"native_api.ListEntitiesCameraResponse":            0x2B,
	"native_api.CameraImageResponse":                   0x2C,
	"native_api.CameraImageRequest":                    0x2D,
	"native_api.ListEntitiesClimateResponse":           0x2E,
	"native_api.ClimateStateResponse":                  0x2F,
	"native_api.ClimateCommandRequest":                 0x30,
}

var messageTypesByID map[byte]string

func init() {
	messageTypesByID = make(map[byte]string, len(messageTypesByName))

	for k, v := range messageTypesByName {
		messageTypesByID[v] = k
	}
}

type handler func(proto.Message, error) bool
type handlerSubscribe func(proto.Message) bool

type Client struct {
	address  string
	password string
	clientID string
	debug    uint32
	mutex    sync.RWMutex

	inited        bool
	initMutex     sync.Mutex
	authenticated uint32
	restarted     uint32

	closeDeadline int64
	conn          net.Conn
	handlers      []handler
	done          chan struct{}
}

func New(address, password string) *Client {
	if _, _, err := net.SplitHostPort(address); err != nil {
		address = address + ":" + strconv.Itoa(defaultPort)
	}

	return &Client{
		address:  address,
		password: password,
		clientID: defaultClient,
		handlers: make([]handler, 0),
		done:     make(chan struct{}),
	}
}

func (c *Client) ClientID() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.clientID
}

func (c *Client) WithClientID(id string) *Client {
	c.mutex.Lock()
	c.clientID = id
	c.mutex.Unlock()

	return c
}

func (c *Client) Debug() bool {
	return atomic.LoadUint32(&c.debug) != 0
}

func (c *Client) WithDebug(debug bool) *Client {
	if debug {
		atomic.StoreUint32(&c.debug, 1)
	} else {
		atomic.StoreUint32(&c.debug, 0)
	}

	return c
}

func (c *Client) Close() (err error) {
	if e := c.connectionCheck(); e == nil {
		err = c.write(context.Background(), &DisconnectRequest{})
	}

	if err != nil {
		return err
	}

	close(c.done)

	c.mutex.Lock()
	c.handlers = c.handlers[:0]
	c.done = make(chan struct{})
	c.mutex.Unlock()

	c.connectionClose()
	c.authenticateClose()

	c.initMutex.Lock()
	c.inited = false
	c.initMutex.Unlock()

	atomic.StoreUint32(&c.restarted, 0)

	return nil
}

func (c *Client) restart() {
	atomic.StoreUint32(&c.restarted, 1)

	c.connectionClose()
	c.authenticateClose()

	c.initMutex.Lock()
	c.inited = false
	c.initMutex.Unlock()
}

func (c *Client) isRestart() bool {
	return atomic.LoadUint32(&c.restarted) == 1
}

func (c *Client) init() (err error) {
	c.initMutex.Lock()
	defer c.initMutex.Unlock()

	if !c.inited {
		now := time.Now()
		deadline := time.Unix(0, atomic.LoadInt64(&c.closeDeadline))

		if deadline.After(now) {
			time.Sleep(deadline.Sub(now))
		}

		err = c.connectionInit()
		if err != nil {
			return err
		}

		c.connectionRun()

		err = c.write(context.Background(), &HelloRequest{ClientInfo: c.ClientID()})
		if err != nil {
			c.connectionClose()
		}

		err = c.authenticateRun()
		if err != nil {
			c.authenticateClose()
		}

		c.inited = true
		atomic.StoreUint32(&c.restarted, 0)
	}

	return err
}

func (c *Client) handlerRegister(h handler) {
	c.mutex.Lock()
	c.handlers = append(c.handlers, h)
	c.mutex.Unlock()
}

func (c *Client) handle(message proto.Message, err error) {
	c.mutex.RLock()
	handlers := make([]handler, len(c.handlers))
	copy(handlers, c.handlers)
	c.mutex.RUnlock()

	for i := len(handlers) - 1; i >= 0; i-- {
		if handlers[i](message, err) {
			handlers = append(handlers[:i], handlers[i+1:]...)
		}
	}

	c.mutex.Lock()
	c.handlers = handlers
	c.mutex.Unlock()
}

func (c *Client) invoke(ctx context.Context, request proto.Message, messageName string) (message proto.Message, err error) {
	chMessage := make(chan proto.Message, 1)
	chErr := make(chan error, 1)

	err = c.invokeHandler(ctx, request, func(message proto.Message, err error) bool {
		if err != nil {
			chErr <- err
			return true
		}

		if messageName == "" || strings.Compare(proto.MessageName(message), messageName) == 0 {
			chMessage <- message
			return true
		}

		return false
	})

	if err != nil {
		return nil, err
	}

	select {
	case message := <-chMessage:
		return message, nil

	case err := <-chErr:
		return nil, err
	}
}

func (c *Client) invokeHandler(ctx context.Context, request proto.Message, h handler) (err error) {
	if err = c.init(); err != nil {
		return err
	}

	if err = c.write(ctx, request); err != nil {
		return err
	}

	c.handlerRegister(h)
	return nil
}

func (c *Client) invokeNoDelay(ctx context.Context, request proto.Message) error {
	if err := c.init(); err != nil {
		return err
	}

	return c.write(ctx, request)
}

func (c *Client) subscribe(ctx context.Context, request proto.Message, h handlerSubscribe) (_ <-chan proto.Message, err error) {
	chMessages := make(chan proto.Message, 1)
	var canceled uint32

	ctx, cancel := context.WithCancel(ctx)

	if err != nil {
		return nil, err
	}

	err = c.invokeHandler(ctx, request, func(message proto.Message, err error) bool {
		if atomic.LoadUint32(&canceled) != 0 {
			close(chMessages)
			return true
		}

		if err != nil {
			close(chMessages)
			return true
		}

		if h(message) {
			chMessages <- message
		}

		return false
	})

	go func() {
		select {
		case <-ctx.Done():

			atomic.StoreUint32(&canceled, 1)
			cancel()

			return
		}
	}()

	return chMessages, err
}
