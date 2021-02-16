package rpi

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/rpi"
)

type Bind struct {
	di.ConfigBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.WorkersBind

	config       *Config
	serialNumber string

	providerVCGenCMD *rpi.VCGenCMD
	providerSysFS    *rpi.SysFS
}

func (b *Bind) Run() error {
	b.Meta().SetSerialNumber(b.serialNumber)

	return nil
}
