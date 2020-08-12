package v3

import (
	"github.com/kihamo/boggart/components/boggart/di"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config   *Config
	provider *mercury.MercuryV3
}
