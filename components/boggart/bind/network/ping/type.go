package ping

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicOnline = config.TopicOnline.Format(config.Hostname)
	config.TopicLatency = config.TopicLatency.Format(config.Hostname)
	config.TopicCheck = config.TopicCheck.Format(config.Hostname)

	return &Bind{
		config: config,
	}, nil
}
