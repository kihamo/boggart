package miio

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config *Config
	device *vacuum.Device
}

func (b *Bind) Close() error {
	return b.device.Close()
}
