package ping

import (
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind

	config *Config
}
