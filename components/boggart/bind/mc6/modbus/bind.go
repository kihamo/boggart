package modbus

import (
	"context"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/modbus"
	"github.com/kihamo/boggart/providers/mc6"
	"go.uber.org/multierr"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	provider     *mc6.MC6
	providerOnce *atomic.Once

	stateDeviceType *atomic.Uint32Null
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

func (b *Bind) DeviceType(ctx context.Context) (mc6.Device, error) {
	if !b.stateDeviceType.IsNil() {
		return mc6.NewDevice(uint16(b.stateDeviceType.Load())), nil
	}

	deviceType, err := b.Provider().DeviceType()

	if err == nil {
		b.stateDeviceType.Set(uint32(deviceType))

		if e := b.MQTT().PublishAsync(ctx, b.config().TopicDeviceType.Format(b.Meta().ID()), deviceType); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return deviceType, err
}

func (b *Bind) notifyChangeState(ctx context.Context) error {
	// TODO: 10 секунд задержки на переключение реле

	return b.Workers().TaskRunByName(ctx, "state-updater")
}
