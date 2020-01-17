package herospeed

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/herospeed"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind

	config *Config
	client *herospeed.Client
}
