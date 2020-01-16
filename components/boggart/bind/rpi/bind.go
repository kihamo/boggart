package rpi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/rpi"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config

	providerVCGenCMD *rpi.VCGenCMD
	providerSysFS    *rpi.SysFS
}
