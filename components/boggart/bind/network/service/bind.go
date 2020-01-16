package service

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config  *Config
	address string
}
