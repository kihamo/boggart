package v3

import (
	"encoding/hex"
	"fmt"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config       *Config
	provider     *mercury.MercuryV3
	providerOnce *atomic.Once
	connection   connection.Connection
}

func (b *Bind) Run() error {
	b.providerOnce.Reset()

	return nil
}

func (b *Bind) Provider() (provider *mercury.MercuryV3, err error) {
	b.providerOnce.Do(func() {
		var conn connection.Connection

		conn, err = connection.NewByDSNString(b.config.ConnectionDSN)
		if err != nil {
			return
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

			b.Logger().Debug("Read packet", args...)
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

			b.Logger().Debug("Write packet", args...)
		}))

		opts := make([]mercury.Option, 0)
		if b.config.Address != "" {
			opts = append(opts, mercury.WithAddressAsString(b.config.Address))
		}

		b.connection = conn
		b.provider = mercury.New(conn, opts...)
	})

	if err != nil {
		b.providerOnce.Reset()
		return nil, err
	}

	return b.provider, err
}

func (b *Bind) Close() error {
	if b.providerOnce.IsDone() {
		return b.connection.Close()
	}

	return nil
}
