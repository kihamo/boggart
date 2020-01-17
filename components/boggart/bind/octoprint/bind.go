package octoprint

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/octoprint"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind

	config   *Config
	provider *octoprint.Client
}
