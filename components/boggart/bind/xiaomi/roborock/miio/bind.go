package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config
	device *vacuum.Device
}

func (b *Bind) Close() error {
	return b.device.Close()
}
