package mqtt

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	if config.TopicLog == "" {
		config.TopicLog = config.TopicPrefix + "/debug"
	}

	if config.TopicBirth == "" {
		config.TopicBirth = config.TopicPrefix + "/status"
	}

	if config.TopicWill == "" {
		config.TopicWill = config.TopicPrefix + "/status"
	}

	bind := &Bind{
		config: config,
	}

	return bind, nil
}
