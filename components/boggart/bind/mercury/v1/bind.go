package v1

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config   *Config
	provider *mercury.MercuryV1
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Address)

	conn, err := connection.NewByDSNString(b.config.ConnectionDSN)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(b.config.Location)
	if err != nil {
		return err
	}

	dump := func(message string) func([]byte) {
		return func(data []byte) {
			args := make([]interface{}, 0)

			packet := mercury.NewPacket()
			if err := packet.UnmarshalBinary(data); err == nil {
				args = append(args,
					"address", fmt.Sprintf("0x%X", packet.Address()),
					"command", fmt.Sprintf("0x%X", packet.Command()),
					"payload", fmt.Sprintf("%v", packet.Payload()),
					"crc", fmt.Sprintf("0x%X", packet.CRC()),
				)
			} else {
				args = append(args,
					"payload", fmt.Sprintf("%v", data),
					"hex", "0x"+hex.EncodeToString(data),
				)
			}

			b.Logger().Debug(message, args...)
		}
	}

	conn.ApplyOptions(connection.WithDumpRead(dump("Read packet")))
	conn.ApplyOptions(connection.WithDumpWrite(dump("Write packet")))

	t := b.Meta().BindType().(Type)

	opts := []mercury.Option{
		t.SerialNumberFunc(b.config.Address),
		mercury.WithDevice(t.Device),
		mercury.WithLocation(loc),
	}

	b.provider = mercury.New(conn, opts...)

	return nil
}
