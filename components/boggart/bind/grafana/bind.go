package grafana

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/grafana"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind

	config   *Config
	provider *grafana.Client
}
