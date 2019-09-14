package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicValue = config.TopicValue.Format(config.Address)

	bind := &Bind{
		config: config,
	}

	bind.SetSerialNumber(config.Address)

	return bind, nil
}
