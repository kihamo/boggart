package pulsar

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/kihamo/boggart/atomic"
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
	di.WidgetBind

	config         *Config
	provider       *pulsar.HeatMeter
	providerMutex  sync.RWMutex
	connection     connection.Connection
	connectionOnce *atomic.Once
	location       *time.Location
}

func (b *Bind) Run() error {
	if b.config.Address != "" {
		if address, err := hex.DecodeString(b.config.Address); err == nil {
			return b.createProvider(address)
		}
	}

	b.connectionOnce.Reset()

	return nil
}

func (b *Bind) getConnection() (conn connection.Connection, err error) {
	b.connectionOnce.Do(func() {
		conn, err = connection.NewByDSNString(b.config.ConnectionDSN)
		if err != nil {
			return
		}

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

		conn.ApplyOptions(connection.WithDumpRead(dump("Read packet")))
		conn.ApplyOptions(connection.WithDumpWrite(dump("Write packet")))

		b.connection = conn
	})

	if err != nil {
		b.connectionOnce.Reset()
		return nil, err
	}

	return b.connection, err
}

func (b *Bind) Provider() *pulsar.HeatMeter {
	b.providerMutex.RLock()
	defer b.providerMutex.RUnlock()

	return b.provider
}

func (b *Bind) createProvider(address []byte) error {
	conn, err := b.getConnection()
	if err != nil {
		return err
	}

	opts := []pulsar.Option{
		pulsar.WithAddress(address),
		pulsar.WithLocation(b.location),
	}

	b.Meta().SetSerialNumber(hex.EncodeToString(address))

	b.providerMutex.Lock()
	b.provider = pulsar.New(conn, opts...)
	b.providerMutex.Unlock()

	return nil
}

func (b *Bind) inputVolume(pulses float32, offset float32) float32 {
	return (offset*InputScale + pulses*10) / InputScale
}

func (b *Bind) Close() error {
	if b.connectionOnce.IsDone() {
		return b.connection.Close()
	}

	return nil
}
