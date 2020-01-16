package octoprint

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/octoprint"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config   *Config
	provider *octoprint.Client
}
