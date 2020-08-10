package v3

import (
	"encoding/hex"
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	conn, err := connection.NewByDSNString(config.ConnectionDSN)
	if err != nil {
		return nil, err
	}

	opts := make([]mercury.Option, 0)
	if config.Address != "" {
		opts = append(opts, mercury.WithAddressAsString(config.Address))
	}

	bind := &Bind{
		provider: mercury.New(conn, opts...),
		config:   config,
	}

	conn.ApplyOptions(connection.WithDumpRead(func(bytes []byte) {
		bind.Logger().Debug("Read packet",
			"payload", fmt.Sprintf("%v", bytes),
			"hex", hex.EncodeToString(bytes),
		)
	}))
	conn.ApplyOptions(connection.WithDumpWrite(func(bytes []byte) {
		bind.Logger().Debug("Write packet",
			"payload", fmt.Sprintf("%v", bytes),
			"hex", hex.EncodeToString(bytes),
		)
	}))

	return bind, nil
}
