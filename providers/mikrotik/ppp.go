package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/types"
)

type PPPActive struct {
	Id            string   `mapstructure:".id"`
	Name          string   `mapstructure:"name"`
	Service       string   `mapstructure:"service"`
	CallerID      string   `mapstructure:"caller-id"`
	Address       types.IP `mapstructure:"address"`
	Uptime        string   `mapstructure:"uptime"`
	Encoding      string   `mapstructure:"encoding"`
	SessionID     string   `mapstructure:"session-id"`
	LimitBytesIn  uint64   `mapstructure:"limit-bytes-in"`
	LimitBytesOut uint64   `mapstructure:"limit-bytes-out"`
	Radius        bool     `mapstructure:"radius"`
	Comment       string   `mapstructure:"comment,omitempty"`
}

func (c *Client) PPPActive(ctx context.Context) (result []PPPActive, err error) {
	err = c.doConvert(ctx, []string{"/ppp/active/print"}, &result)
	return result, err
}
