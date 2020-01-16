package grafana

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/grafana"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind

	config   *Config
	provider *grafana.Client
}
