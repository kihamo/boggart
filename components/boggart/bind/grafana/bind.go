package grafana

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/grafana"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *grafana.Client
}
