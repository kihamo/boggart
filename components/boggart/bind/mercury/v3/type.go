package v3

import (
	"encoding/hex"
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Type struct{}

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

	conn.ApplyOptions(connection.WithDumpRead(func(data []byte) {
		args := make([]interface{}, 0)

		packet := mercury.NewResponse()
		if err := packet.UnmarshalBinary(data); err == nil {
			args = append(args,
				"address", fmt.Sprintf("0x%X", packet.Address()),
				"payload", fmt.Sprintf("%v", packet.Payload()),
				"crc", fmt.Sprintf("0x%X", packet.CRC()),
			)
		} else {
			args = append(args,
				"payload", fmt.Sprintf("%v", data),
				"hex", "0x"+hex.EncodeToString(data),
			)
		}

		bind.Logger().Debug("Read packet", args...)
	}))
	conn.ApplyOptions(connection.WithDumpWrite(func(data []byte) {
		args := make([]interface{}, 0)

		packet := mercury.NewRequest()
		if err := packet.UnmarshalBinary(data); err == nil {
			args = append(args,
				"address", fmt.Sprintf("0x%X", packet.Address()),
				"code", fmt.Sprintf("0x%X", packet.Code()),
				"parameter.code", fmt.Sprintf("0x%X", packet.ParameterCode()),
				"parameter.extension", fmt.Sprintf("0x%X", packet.ParameterExtension()),
				"parameters", fmt.Sprintf("%v", packet.Parameters()),
				"crc", fmt.Sprintf("0x%X", packet.CRC()),
			)
		} else {
			args = append(args,
				"payload", fmt.Sprintf("%v", data),
				"hex", "0x"+hex.EncodeToString(data),
			)
		}

		bind.Logger().Debug("Write packet", args...)
	}))

	return bind, nil
}
