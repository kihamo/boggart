package tvt

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/tvt"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config *Config
	client *tvt.Client
}
