package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/rkcm"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config
	client *rkcm.Client
}
