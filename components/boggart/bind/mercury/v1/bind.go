package v1

import (
	"github.com/kihamo/boggart/components/boggart/di"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config   *Config
	provider *mercury.MercuryV1
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Address)

	return nil
}
