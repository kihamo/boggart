package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicName = config.TopicName.Format(config.ApplicationName)
	config.TopicVersion = config.TopicVersion.Format(config.ApplicationName)
	config.TopicBuild = config.TopicBuild.Format(config.ApplicationName)
	config.TopicShutdown = config.TopicShutdown.Format(config.ApplicationName)

	bind := &Bind{
		config: config,
	}

	return bind, nil
}
