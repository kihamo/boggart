package mikrotik

import (
	"context"
)

type PPPActive struct {
	Id            string `json:".id"`
	Name          string `json:"name"`
	Service       string `json:"service"`
	CallerID      string `json:"caller-id"`
	Address       string `json:"address"`
	Uptime        string `json:"uptime"`
	Encoding      string `json:"encoding"`
	SessionID     string `json:"session-id"`
	LimitBytesIn  uint64 `json:"limit-bytes-in"`
	LimitBytesOut uint64 `json:"limit-bytes-out"`
	Radius        bool   `json:"radius"`
	Comment       string `json:"comment,omitempty"`
}

func (c *Client) PPPActive(ctx context.Context) (result []PPPActive, err error) {
	err = c.doConvert(ctx, []string{"/ppp/active/print"}, &result)
	return result, err
}
