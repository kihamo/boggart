package rpi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/rpi"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config

	providerVCGenCMD *rpi.VCGenCMD
	providerSysFS    *rpi.SysFS
}
