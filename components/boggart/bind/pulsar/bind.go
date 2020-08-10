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
	if err != nil {
		b.connection.ApplyOptions(connection.WithDumpRead(func(bytes []byte) {
			b.Logger().Debug("Read packet",
				"payload", fmt.Sprintf("%v", bytes),
				"hex", hex.EncodeToString(bytes),
			)
		}))
		b.connection.ApplyOptions(connection.WithDumpWrite(func(bytes []byte) {
			b.Logger().Debug("Write packet",
				"payload", fmt.Sprintf("%v", bytes),
				"hex", hex.EncodeToString(bytes),
			)
		}))
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
