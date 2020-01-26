package miio

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config *Config
	device *vacuum.Device
}

func (b *Bind) Close() error {
	return b.device.Close()
}
