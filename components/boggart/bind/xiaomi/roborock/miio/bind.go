package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
	device *vacuum.Device

	state *atomic.Uint32Null

	battery   *atomic.Uint32Null
	cleanArea *atomic.Uint32Null
	cleanTime *atomic.DurationNull
	fanPower  *atomic.Uint32Null
	volume    *atomic.Uint32Null

	consumableFilter    *atomic.DurationNull
	consumableBrushMain *atomic.DurationNull
	consumableBrushSide *atomic.DurationNull
	consumableSensor    *atomic.DurationNull
}

func (b *Bind) Close() error {
	return b.device.Close()
}
