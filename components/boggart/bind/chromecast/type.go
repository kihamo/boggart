package chromecast

import (
	"net"
	"strconv"

	"github.com/barnybug/go-cast/events"
	"github.com/barnybug/go-cast/log"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	sn := net.JoinHostPort(config.Host.String(), strconv.Itoa(config.Port))

	config.TopicVolume = config.TopicVolume.Format(sn)
	config.TopicMute = config.TopicMute.Format(sn)
	config.TopicPause = config.TopicPause.Format(sn)
	config.TopicStop = config.TopicStop.Format(sn)
	config.TopicPlay = config.TopicPlay.Format(sn)
	config.TopicResume = config.TopicResume.Format(sn)
	config.TopicSeek = config.TopicSeek.Format(sn)
	config.TopicAction = config.TopicAction.Format(sn)
	config.TopicStateStatus = config.TopicStateStatus.Format(sn)
	config.TopicStateVolume = config.TopicStateVolume.Format(sn)
	config.TopicStateMute = config.TopicStateMute.Format(sn)
	config.TopicStateContent = config.TopicStateContent.Format(sn)

	log.Debug = config.Debug

	bind := &Bind{
		config:         config,
		disconnected:   atomic.NewBoolNull(),
		volume:         atomic.NewUint32Null(),
		mute:           atomic.NewBoolNull(),
		status:         atomic.NewString(),
		mediaContentID: atomic.NewString(),
		events:         make(chan events.Event, 16),
	}
	bind.SetSerialNumber(sn)

	return bind, nil
}
