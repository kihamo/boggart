package v1

import (
	"github.com/kihamo/boggart/components/boggart/di"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind

	config   *Config
	provider *mercury.MercuryV1
}
