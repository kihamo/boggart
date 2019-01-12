package lg_webos

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	device := &Bind{
		config: c.(*Config),
	}
	device.Init()

	return device, nil
}
