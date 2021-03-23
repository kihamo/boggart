package ledwifi

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/wifiled"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	bulb *wifiled.Bulb
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.bulb = wifiled.NewBulb(b.config().Address)

	return nil
}
