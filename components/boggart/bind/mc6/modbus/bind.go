package modbus

import (
	"context"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/modbus"
	"github.com/kihamo/boggart/providers/mc6"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	provider     *mc6.MC6
	providerOnce *atomic.Once

	stateDeviceType *atomic.Uint32
	statePower      *atomic.Bool
	//stateTemperatureFormat *atomic.Uint32
	stateSetTemperature     *atomic.Float64
	stateAway               *atomic.Bool
	stateAwayTemperature    *atomic.Uint32
	stateHoldingTemperature *atomic.Uint32
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.providerOnce.Reset()

	return nil
}

func (b *Bind) Provider() (provider *mc6.MC6) {
	b.providerOnce.Do(func() {
		cfg := b.config()

		b.provider = mc6.New(
			&b.config().DSN.URL,
			mc6.WithSlaveID(cfg.ConnectionSlaveID),
			mc6.WithTimeout(cfg.ConnectionTimeout),
			mc6.WithIdleTimeout(cfg.ConnectionIdleTimeout),
			mc6.WithLogger(modbus.NewLogger(func(s string) {
				b.Logger().Debug(s)
			})),
		)
	})

	return b.provider
}

func (b *Bind) Close() error {
	if b.providerOnce.IsDone() {
		return b.provider.Close()
	}

	return nil
}

func (b *Bind) Power(ctx context.Context, flag bool) error {
	err := b.Provider().Status(flag)

	if err == nil {
		b.statePower.Set(flag)

		err = b.MQTT().PublishAsync(ctx, b.config().TopicPowerState.Format(b.Meta().ID()), flag)
	}

	return err
}

//func (b *Bind) TemperatureFormat(ctx context.Context, format uint16) error {
//	err := b.Provider().TemperatureFormat(format)
//
//	if err == nil {
//		b.stateTemperatureFormat.Set(uint32(format))
//
//		err = b.MQTT().PublishAsync(ctx, b.config().TopicTemperatureFormatState.Format(b.Meta().ID()), format)
//	}
//
//	return err
//}

func (b *Bind) SetTemperature(ctx context.Context, temperature float64) error {
	err := b.Provider().SetTemperature(temperature)

	if err == nil {
		// устанавливаемое значение всегда кратно 0.5 и округляется в меньшую сторону
		// даже на устройстве шаг 0.5, поэтому принудительно округляем
		val := int(temperature * 10)
		val -= val % 5
		temperature = float64(val) / 10

		b.stateSetTemperature.Set(temperature)

		err = b.MQTT().PublishAsync(ctx, b.config().TopicSetTemperatureState.Format(b.Meta().ID()), temperature)
	}

	return err
}

func (b *Bind) Away(ctx context.Context, flag bool) error {
	err := b.Provider().Away(flag)

	if err == nil {
		b.stateAway.Set(flag)

		err = b.MQTT().PublishAsync(ctx, b.config().TopicAwayState.Format(b.Meta().ID()), flag)
	}

	return err
}

func (b *Bind) AwayTemperature(ctx context.Context, temperature uint16) error {
	err := b.Provider().AwayTemperature(temperature)

	if err == nil {
		b.stateAwayTemperature.Set(uint32(temperature))

		err = b.MQTT().PublishAsync(ctx, b.config().TopicAwayTemperatureState.Format(b.Meta().ID()), temperature)
	}

	return err
}

func (b *Bind) HoldingTemperature(ctx context.Context, temperature uint16) error {
	err := b.Provider().HoldingTemperature(temperature)

	if err == nil {
		b.stateHoldingTemperature.Set(uint32(temperature))

		err = b.MQTT().PublishAsync(ctx, b.config().TopicHoldingTemperatureState.Format(b.Meta().ID()), temperature)
	}

	return err
}
