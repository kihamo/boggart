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

	providerVCGenCMD *rpi.VCGenCMD
	providerSysFS    *rpi.SysFS
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.providerSysFS = rpi.NewSysFS()
	b.providerVCGenCMD = rpi.NewVCGenCMD()

	sn, err := b.providerSysFS.SerialNumber()
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(sn)

	return nil
}
