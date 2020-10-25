package rkcm

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/rkcm"
)

type Bind struct {
	di.LoggerBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	config *Config
	client *rkcm.Client
}
