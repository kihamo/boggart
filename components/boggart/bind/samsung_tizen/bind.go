package tizen

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	client *tv.APIv2
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	if cfg.MAC != nil && b.Meta().MAC() == nil {
		b.Meta().SetMAC(cfg.MAC.HardwareAddr)
	}

	b.client = tv.NewAPIv2(cfg.Host)

	return nil
}
