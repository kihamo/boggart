package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config
	client *elektroset.Client
}
