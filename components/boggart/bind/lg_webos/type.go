package webos

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	if config.MAC != nil {
		mac := config.MAC.String()

		config.TopicApplication = config.TopicApplication.Format(mac)
		config.TopicMute = config.TopicMute.Format(mac)
		config.TopicVolume = config.TopicVolume.Format(mac)
		config.TopicVolumeUp = config.TopicVolumeUp.Format(mac)
		config.TopicVolumeDown = config.TopicVolumeDown.Format(mac)
		config.TopicToast = config.TopicToast.Format(mac)
		config.TopicPower = config.TopicPower.Format(mac)
		config.TopicStateMute = config.TopicStateMute.Format(mac)
		config.TopicStateVolume = config.TopicStateVolume.Format(mac)
		config.TopicStateApplication = config.TopicStateApplication.Format(mac)
		config.TopicStateChannelNumber = config.TopicStateChannelNumber.Format(mac)
		config.TopicStatePower = config.TopicStatePower.Format(mac)
	}

	bind := &Bind{
		config: config,
		power:  atomic.NewBoolNull(),
	}

	return bind, nil
}
