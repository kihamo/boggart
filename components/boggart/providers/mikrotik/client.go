package mikrotik

import (
	"context"
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

func (c *Client) doConvert(ctx context.Context, sentence []string, result interface{}) error {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "call")
	defer span.Finish()

	span.SetTag("sentence", strings.Join(sentence, " "))

	c.mutex.Lock()
	reply, err := c.client.RunArgs(sentence)
	c.mutex.Unlock()

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
