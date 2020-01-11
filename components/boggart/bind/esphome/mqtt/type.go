package mqtt

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	if config.TopicLog == "" {
		config.TopicLog = config.TopicPrefix + "/debug"
	}

	bind := &Bind{
		config: config,
	}

	return bind, nil
}
