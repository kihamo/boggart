package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	conn, err := connection.New(config.ConnectionDSN)
	if err != nil {
		return nil, err
	}

	opts := make([]mercury.Option, 0)
	if config.Address != "" {
		opts = append(opts, mercury.WithAddressAsString(config.Address))
	}

	return &Bind{
		provider: mercury.New(conn, opts...),
		config:   config,
	}, nil
}
