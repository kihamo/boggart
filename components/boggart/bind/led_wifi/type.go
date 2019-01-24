package led_wifi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/wifiled"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device := &Bind{
		bulb:  wifiled.NewBulb(config.Address),
		power: atomic.NewBoolNull(),
		mode:  atomic.NewUint32Null(),
		speed: atomic.NewUint32Null(),
		color: atomic.NewUint32Null(),
	}

	return device, nil
}
