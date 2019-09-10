package mikrotik

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow/components/tracing"
	"gopkg.in/routeros.v2"
)

const (
	ComponentName = "mikrotik"
)

var (
	ErrEmptyResponse = errors.New("empty response")
)

func IsEmptyResponse(err error) bool {
	return err == ErrEmptyResponse
}

type Client struct {
	address  string
	username string
	password string

	dialer *net.Dialer
	conn   net.Conn

	client    *routeros.Client
	connected uint64
	mutex     sync.Mutex
}

func NewClient(address, username, password string, timeout time.Duration) *Client {
	return &Client{
		address:  address,
		username: username,
		password: password,
		dialer: &net.Dialer{
			Timeout: timeout,
		},
	}
}

func (c *Client) connect() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn, err = c.dialer.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	if c.dialer.Timeout > 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.dialer.Timeout)); err != nil {
			return err
		}
	}

	c.client, err = routeros.NewClient(c.conn)
	if err != nil {
		c.conn.Close()
		return err
	}

	err = c.client.Login(c.username, c.password)
	if err != nil {
		c.conn.Close()
		return err
	}

	atomic.StoreUint64(&c.connected, 1)

	return nil
}

func (c *Client) isRetry(err error) bool {
	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		return true
	}

	return err == io.EOF
}

func (c *Client) safeRunArgs(sentence []string) (*routeros.Reply, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.client.RunArgs(sentence)
}

func (c *Client) doConvert(ctx context.Context, sentence []string, result interface{}) error {
	if atomic.LoadUint64(&c.connected) == 0 {
		if err := c.connect(); err != nil {
			return err
		}
	}

	if c.dialer.Timeout > 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.dialer.Timeout)); err != nil {
			return err
		}
	}

	span, _ := tracing.StartSpanFromContext(ctx, ComponentName, "call")
	defer span.Finish()

	span.SetTag("sentence", strings.Join(sentence, " "))

	reply, err := c.safeRunArgs(sentence)
	if err != nil && c.isRetry(err) {
		atomic.StoreUint64(&c.connected, 0)

		if err = c.connect(); err == nil {
			reply, err = c.safeRunArgs(sentence)
		}
	}

	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	if len(reply.Re) == 0 || (len(reply.Re[0].List) == 0 && len(reply.Re[0].Map) == 0) {
		return nil
	}

	records := make([]interface{}, 0, len(reply.Re))
	for _, re := range reply.Re {
		records = append(records, re.Map)
	}

	converter := gotypes.NewConverter(records, result)
	if !converter.Valid() {
		err = fmt.Errorf("failed convert fields: %v", strings.Join(converter.GetInvalidFields(), ","))
	}

	if err != nil {
		tracing.SpanError(span, err)
	}

	return err
}
