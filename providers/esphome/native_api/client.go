package native_api

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	defaultClient = "GoLang Native API client"

	packetMagicByte   byte = 0x00
	keepAliveInterval      = time.Second * 5

	authNone int32 = iota
	authSuccess
	authFailed
)

var (
	ErrAuthenticated      = errors.New("must authenticated")
	ErrAuthenticateFailed = errors.New("authenticated failed")
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

type Client struct {
	address       string
	password      string
	clientID      string
	debug         uint32
	mutex         sync.RWMutex
	done          chan struct{}
	authenticated int32

	connection     *connection
	connectionID   uint64
	connectionPing bool
}

func New(address, password string) *Client {
	return &Client{
		address:        address,
		password:       password,
		clientID:       defaultClient,
		connectionID:   1,
		connectionPing: true,
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

	c.mutex.RLock()
	if c.connection != nil {
		c.connection.WithDebug(debug)
	}
	c.mutex.RUnlock()

	return c
}

func (c *Client) connect() (conn *connection, err error) {
	c.mutex.RLock()
	conn = c.connection
	address := c.address
	c.mutex.RUnlock()

	if conn == nil {
		conn, err = newConnection(address)
		if err != nil {
			return nil, err
		}

		conn.WithDebug(c.Debug())
		conn.WithID(c.connectionID)

		c.mutex.Lock()
		c.connection = conn
		c.mutex.Unlock()

		ctx := context.Background()

		_, err = c.Hello(ctx)
		if err != nil {
			c.mutex.Lock()
			c.connection = nil
			c.mutex.Unlock()

			return nil, err
		}

		if c.connectionPing {
			go c.keepalive()
		}
	}

	return conn, err
}

func (c *Client) Clone() (*Client, error) {
	c.mutex.RLock()
	clientId := c.clientID
	c.mutex.RUnlock()

	client := New(c.address, c.password).
		WithDebug(c.Debug()).
		WithClientID(clientId)
	client.connectionID = c.connectionID + 1

	return client, nil
}

func (c *Client) Close() error {
	c.mutex.RLock()
	conn := c.connection
	c.mutex.RUnlock()

	if conn != nil {
		if _, err := c.invoke(context.Background(), &DisconnectRequest{}); err != nil {
			return err
		}

		if err := conn.Close(); err != nil {
			return err
		}

		c.mutex.Lock()
		c.connection = nil
		close(c.done)
		c.mutex.Unlock()
	}

	atomic.StoreInt32(&c.authenticated, authNone)

	return nil
}

func (c *Client) invoke(ctx context.Context, request proto.Message) (proto.Message, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	conn.Lock()
	defer conn.Unlock()

	if err := conn.Write(ctx, request); err != nil {
		return nil, err
	}

	return conn.Read(ctx)
}

func (c *Client) invokeNoDelay(ctx context.Context, request proto.Message) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}

	conn.Lock()
	defer conn.Unlock()

	return conn.Write(ctx, request)
}

func (c *Client) keepalive() {
	done := make(chan struct{})

	c.mutex.Lock()
	c.done = done
	c.mutex.Unlock()

	ticker := time.NewTicker(keepAliveInterval)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			_, err := c.Ping(context.Background())
			if err != nil {
				c.Close()
				return
			}

		case <-done:
			return
		}
	}
}

func (c *Client) checkAuthenticate(ctx context.Context) error {
	switch atomic.LoadInt32(&c.authenticated) {
	case authSuccess:
		return nil

	case authFailed:
		return ErrAuthenticated

	default:
		response, err := c.Connect(ctx)
		if err != nil {
			atomic.StoreInt32(&c.authenticated, authNone)
			return err
		}

		if response.GetInvalidPassword() {
			atomic.StoreInt32(&c.authenticated, authFailed)
			return ErrAuthenticateFailed
		}

		atomic.StoreInt32(&c.authenticated, authSuccess)
		return nil
	}

	return nil
}
