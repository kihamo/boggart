package elektroset

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.ProbesBind

	config *Config
	client *elektroset.Client
}
