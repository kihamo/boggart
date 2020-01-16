package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config   *Config
	provider *mercury.MercuryV3
}
