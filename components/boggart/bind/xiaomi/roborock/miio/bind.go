package miio

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
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

	device *vacuum.Device
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()
	b.device = vacuum.New(cfg.Host, cfg.Token)

	if cfg.PacketsCounter > 0 {
		b.device.Client().SetPacketsCounter(cfg.PacketsCounter)
	}

	return nil
}

func (b *Bind) Close() error {
	return b.device.Close()
}
