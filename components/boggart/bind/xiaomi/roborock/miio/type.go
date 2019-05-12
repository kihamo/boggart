package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
		device: devices.NewVacuum(config.Host, config.Token),

		battery:   atomic.NewUint32Null(),
		cleanArea: atomic.NewUint32Null(),
		cleanTime: atomic.NewUint32Null(),
	}

	return bind, nil
}
