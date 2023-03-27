package neptun

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/modbus"
	"github.com/kihamo/boggart/providers/neptun"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.ProbesBind
	di.WidgetBind

	provider     *neptun.Neptun
	providerOnce *atomic.Once
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.providerOnce.Reset()

	return nil
}

func (b *Bind) Provider() (provider *neptun.Neptun) {
	b.providerOnce.Do(func() {
		cfg := b.config()

		b.provider = neptun.New(
			&b.config().DSN.URL,
			modbus.WithSlaveID(cfg.ConnectionSlaveID),
			modbus.WithTimeout(cfg.ConnectionTimeout),
			modbus.WithIdleTimeout(cfg.ConnectionIdleTimeout),
			modbus.WithLogger(modbus.NewLogger(func(s string) {
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
