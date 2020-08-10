package pulsar

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/connection"
	"github.com/kihamo/boggart/providers/pulsar"
)

const (
	InputScale = 1000
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config   *Config
	provider *pulsar.HeatMeter

	location   *time.Location
	connection connection.Connection
}

func (b *Bind) Run() (err error) {
	b.connection, err = connection.NewByDSNString(b.config.ConnectionDSN)
	if err == nil {
		dump := func(message string) func([]byte) {
			return func(data []byte) {
				args := make([]interface{}, 0)

				packet := pulsar.NewPacket()
				if err := packet.UnmarshalBinary(data); err == nil {
					args = append(args,
						"address", "0x"+hex.EncodeToString(packet.Address()),
						"function", fmt.Sprintf("0x%X", packet.Function()),
						"length", packet.Length(),
						"error-code", fmt.Sprintf("0x%X", packet.ErrorCode()),
						"payload", fmt.Sprintf("%v", packet.Payload()),
						"id", "0x"+hex.EncodeToString(packet.ID()),
						"crc", "0x"+hex.EncodeToString(packet.CRC()),
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

		b.connection.ApplyOptions(connection.WithDumpRead(dump("Read packet")))
		b.connection.ApplyOptions(connection.WithDumpWrite(dump("Write packet")))
	}

	if b.config.Address != "" {
		var address []byte

		if address, err = hex.DecodeString(b.config.Address); err == nil {
			b.createProvider(address)
		}
	}

	return err
}

func (b *Bind) createProvider(address []byte) {
	opts := []pulsar.Option{
		pulsar.WithAddress(address),
		pulsar.WithLocation(b.location),
	}

	b.provider = pulsar.New(b.connection, opts...)
	b.Meta().SetSerialNumber(hex.EncodeToString(address))
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}
