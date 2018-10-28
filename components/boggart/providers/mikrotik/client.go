package mikrotik

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/gotypes"
	"github.com/kihamo/shadow/components/tracing"
	"gopkg.in/routeros.v2"
)

const (
	ComponentName = "mikrotik"
)

type Client struct {
	mutex  sync.Mutex
	client *routeros.Client
}

func NewClient(address, username, password string, timeout time.Duration) (*Client, error) {
	client, err := routeros.DialTimeout(address, username, password, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) do(sentence []string) (*routeros.Reply, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.client.RunArgs(sentence)
}

func (c *Client) doConvert(ctx context.Context, sentence []string, result interface{}) error {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "call")
	defer span.Finish()

	span.SetTag("sentence", strings.Join(sentence, " "))

	reply, err := c.do(sentence)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	records := make([]interface{}, 0, len(reply.Re))
	for _, re := range reply.Re {
		records = append(records, re.Map)
	}

	converter := gotypes.NewConverter(records, result)
	if !converter.Valid() {
		err = fmt.Errorf("Failed convert fields: %v", strings.Join(converter.GetInvalidFields(), ","))
	}

	if err != nil {
		tracing.SpanError(span, err)
	}

	return err
}

func (c *Client) WifiClients() ([]map[string]string, error) {
	reply, err := c.do([]string{"/interface/wireless/registration-table/print"})
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, errors.New("Empty reply from device")
	}

	clients := make([]map[string]string, 0, len(reply.Re))

	for _, re := range reply.Re {
		clients = append(clients, re.Map)
	}

	return clients, nil
}

func (c *Client) EthernetStats() ([]map[string]string, error) {
	reply, err := c.do([]string{"/interface/print", "stats"})
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, errors.New("Empty reply from device")
	}

	stats := make([]map[string]string, 0, len(reply.Re))

	for _, re := range reply.Re {
		stats = append(stats, re.Map)
	}

	return stats, nil
}
