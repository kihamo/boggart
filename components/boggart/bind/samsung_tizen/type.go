package samsung_tizen

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	device := &Bind{
		client: tv.NewApiV2(config.Host),
	}
	device.Init()

	return device, nil
}
