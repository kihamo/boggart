package dom24

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/dom24"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	//di.MetricsBind
	//di.MQTTBind
	//di.ProbesBind
	di.WidgetBind

	config   *Config
	provider *dom24.Client
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.config.Phone)

	return nil
}
