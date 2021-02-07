package v1

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config       *Config
	provider     *mercury.MercuryV1
	providerOnce *atomic.Once
	connection   connection.Connection
	tariffCount  *atomic.Uint32Null
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Address)
	b.providerOnce.Reset()
	b.tariffCount.Nil()

	return nil
}

func (b *Bind) TariffCount() (uint32, error) {
	if b.tariffCount.IsNil() {
		provider, err := b.Provider()
		if err != nil {
			return 0, err
		}

		tariffCount, err := provider.TariffCount()
		if err != nil {
			return 0, err
		}

		b.tariffCount.Set(uint32(tariffCount))
	}

	return b.tariffCount.Load(), nil
}

func (b *Bind) Provider() (provider *mercury.MercuryV1, err error) {
	b.providerOnce.Do(func() {
		var (
			conn connection.Connection
			dsn  *connection.DSN
		)

		dsn, err = connection.ParseDSN(b.config.ConnectionDSN)
		if err != nil {
			return
		}

		if dsn.DialTimeout == nil {
			dsn.DialTimeout = &[]time.Duration{time.Second}[0]
		}

		conn, err = connection.NewByDSN(dsn)
		if err != nil {
			return
		}

		loc, err := time.LoadLocation(b.config.Location)
		if err != nil {
			return
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
