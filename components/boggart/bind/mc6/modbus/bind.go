package modbus

import (
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
