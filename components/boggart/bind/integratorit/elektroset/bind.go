package elektroset

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

const (
	layoutPeriod = "2006-01-02"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.ConfigBind
	di.MetaBind
	di.WidgetBind

	config *Config
	client *elektroset.Client
}
