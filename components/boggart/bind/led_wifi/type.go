package led_wifi

import (
	"math"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device := &Bind{
		bulb:       wifiled.NewBulb(config.Address),
		statePower: 0,
		stateMode:  math.MaxUint64,
		stateSpeed: math.MaxUint64,
		stateColor: math.MaxUint64,
	}
	device.Init()

	return device, nil
}
