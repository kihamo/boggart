package alsa

import (
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	sn, err := machineid.ID()
	if err != nil {
		sn, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	config.TopicVolume = config.TopicVolume.Format(sn)
	config.TopicMute = config.TopicMute.Format(sn)
	config.TopicPause = config.TopicPause.Format(sn)
	config.TopicStop = config.TopicStop.Format(sn)
	config.TopicPlay = config.TopicPlay.Format(sn)
	config.TopicResume = config.TopicResume.Format(sn)
	config.TopicAction = config.TopicAction.Format(sn)
	config.TopicStateStatus = config.TopicStateStatus.Format(sn)
	config.TopicStateVolume = config.TopicStateVolume.Format(sn)
	config.TopicStateMute = config.TopicStateMute.Format(sn)

	bind := &Bind{
		config:       config,
		done:         make(chan struct{}, 1),
		playerStatus: atomic.NewInt64Default(StatusStopped.Int64()),
		volume:       atomic.NewInt64Default(config.Volume),
		mute:         atomic.NewBoolDefault(config.Mute),
	}

	return bind, nil
}
