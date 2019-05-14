package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices/vacuum"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		device: vacuum.New(config.Host, config.Token),

		battery:   atomic.NewUint32Null(),
		cleanArea: atomic.NewUint32Null(),
		cleanTime: atomic.NewDurationNull(),
		fanPower:  atomic.NewUint32Null(),
		volume:    atomic.NewUint32Null(),

		consumableFilter:    atomic.NewDurationNull(),
		consumableBrushMain: atomic.NewDurationNull(),
		consumableBrushSide: atomic.NewDurationNull(),
		consumableSensor:    atomic.NewDurationNull(),
	}

	bind.device.Client().SetPacketsCounter(config.PacketsCounter)

	return bind, nil
}
