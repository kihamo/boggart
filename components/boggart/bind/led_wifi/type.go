package ledwifi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/wifiled"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicPower = config.TopicPower.Format(config.Address)
	config.TopicColor = config.TopicColor.Format(config.Address)
	config.TopicMode = config.TopicMode.Format(config.Address)
	config.TopicSpeed = config.TopicSpeed.Format(config.Address)
	config.TopicStatePower = config.TopicStatePower.Format(config.Address)
	config.TopicStateColor = config.TopicStateColor.Format(config.Address)
	config.TopicStateColorHSV = config.TopicStateColorHSV.Format(config.Address)
	config.TopicStateMode = config.TopicStateMode.Format(config.Address)
	config.TopicStateSpeed = config.TopicStateSpeed.Format(config.Address)

	bind := &Bind{
		config: config,
		bulb:   wifiled.NewBulb(config.Address),
	}

	return bind, nil
}
