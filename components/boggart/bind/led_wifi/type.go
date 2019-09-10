package led_wifi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/wifiled"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		bulb: wifiled.NewBulb(config.Address),
	}

	return bind, nil
}
