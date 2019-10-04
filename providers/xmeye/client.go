package xmeye

import (
	"context"
	"io"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	protocol "github.com/kihamo/boggart/protocols/connection"
)

const (
	DefaultTimeout = time.Second
	DefaultPort    = 34567

	defaultPayloadBuffer = 1024

	TimeLayout = "2006-01-02 15:04:05"

	CmdLoginResponse                uint16 = 1000
	CmdLogoutResponse               uint16 = 1002
	CmdKeepAliveResponse            uint16 = 1006
	CmdSystemInfoRequest            uint16 = 1020
	CmdConfigGetRequest             uint16 = 1042
	CmdDefaultConfigGetRequest      uint16 = 1044
	CmdConfigChannelTitleSetRequest uint16 = 1046
	CmdConfigChannelTitleGetRequest uint16 = 1048
	CmdAbilityGetRequest            uint16 = 1360
	CmdPlayRequest                  uint16 = 1420
	CmdPlayClaimRequest             uint16 = 1424
	CmdFileSearchRequest            uint16 = 1440
	CmdLogSearchRequest             uint16 = 1442
	CmdSysManagerRequest            uint16 = 1450
	CmdSysManagerResponse           uint16 = 1451
	CmdTimeRequest                  uint16 = 1452
	CmdDiskManagerRequest           uint16 = 1460
	CmdGuardRequest                 uint16 = 1500
	CmdUnGuardRequest               uint16 = 1502
	CmdAlarmRequest                 uint16 = 1504
	CmdConfigExportRequest          uint16 = 1542
	CmdLogExportRequest             uint16 = 1544
)

type Client struct {
	channelsCount uint64
	extraChannel  uint64
	alarmStarted  uint32

	username string
	password []byte

	dsn        string
	connection *connection

	mutex sync.RWMutex
	done  chan struct{}
}

func New(host, username, password string) (*Client, error) {
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = host + ":" + strconv.Itoa(DefaultPort)
	}

	client := &Client{
		username:     username,
		password:     []byte(password),
		dsn:          "tcp://" + host + "?read-timeout=10s&write-timeout=10s&once=true",
		extraChannel: math.MaxUint64,
	}

	dial, err := protocol.New(client.dsn)
	if err != nil {
		return nil, err
	}

	client.connection = &connection{
		Conn:           dial,
		sessionID:      0,
		sequenceNumber: 0,
	}

	return client, nil
}

func (c *Client) IsAuth() bool {
	return c.connection.SessionID() > 0
}

func (c *Client) Call(ctx context.Context, code uint16, payload interface{}) (*packet, error) {
	if !c.IsAuth() && code != CmdLoginResponse && code != CmdLogoutResponse {
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
	}

	packet := newPacket()
	packet.MessageID = code

	if err := packet.LoadPayload(payload); err != nil {
		return nil, err
	}

	if err := c.connection.send(packet); err != nil {
		return nil, err
	}

	p, err := c.connection.receive()
	if err != nil {
		return nil, err
	}

	return p, p.Payload.Error()
}

func (c *Client) CallWithResult(ctx context.Context, code uint16, payload, result interface{}) error {
	packet, err := c.Call(ctx, code, payload)
	if err != nil {
		return err
	}

	if writer, ok := result.(io.Writer); ok {
		_, err = packet.Payload.WriteTo(writer)
		return err
	}

	return packet.Payload.UnmarshalJSON(result)
}

func (c *Client) Cmd(ctx context.Context, code uint16, cmd string) (*packet, error) {
	return c.Call(ctx, code, map[string]string{
		"Name":      cmd,
		"SessionID": c.connection.SessionIDAsString(),
	})
}

func (c *Client) CmdWithResult(ctx context.Context, code uint16, cmd string, result interface{}) error {
	return c.CallWithResult(ctx, code, map[string]string{
		"Name":      cmd,
		"SessionID": c.connection.SessionIDAsString(),
	}, &result)
}

func (c *Client) Close() (err error) {
	c.mutex.Lock()

	if c.done != nil {
		close(c.done)
		c.done = nil
	}

	if c.connection != nil {
		err = c.connection.Close()
		c.connection = nil
	}

	c.mutex.Unlock()

	return err
}

func (c *Client) ChannelsCount() (count uint64) {
	count = atomic.LoadUint64(&c.channelsCount)

	if count == 0 && !c.IsAuth() {
		if err := c.Login(context.Background()); err == nil {
			count = atomic.LoadUint64(&c.channelsCount)
		}
	}

	return count
}

func (c *Client) ExtraChannel() (channel uint64) {
	channel = atomic.LoadUint64(&c.extraChannel)

	if channel == math.MaxUint64 && !c.IsAuth() {
		if err := c.Login(context.Background()); err == nil {
			channel = atomic.LoadUint64(&c.extraChannel)
		}
	}

	return channel
}
