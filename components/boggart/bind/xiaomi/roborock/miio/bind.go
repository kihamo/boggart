package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
	device *vacuum.Device
}

func (b *Bind) Close() error {
	return b.device.Close()
}
