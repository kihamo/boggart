package rkcm

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/rkcm"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind

	config *Config
	client *rkcm.Client
}
