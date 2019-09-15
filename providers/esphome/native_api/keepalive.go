package native_api

import (
	"context"
	"time"
)

const (
	keepAliveInterval = time.Second * 5
)

func (c *Client) keepalive() {
	ticker := time.NewTicker(keepAliveInterval)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-c.done:
			return

		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), defaultWriteTimeout)
			c.Ping(ctx)
			cancel()
		}
	}
}
