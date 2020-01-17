package grafana

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/grafana"
)

type Bind struct {
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config   *Config
	provider *grafana.Client
}
