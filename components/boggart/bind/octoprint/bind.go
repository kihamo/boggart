package octoprint

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/octoprint"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *octoprint.Client
}
