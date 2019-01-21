package alsa

import (
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	sn, err := machineid.ID()
	if err != nil {
		sn, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	device := &Bind{
		done:         make(chan struct{}, 1),
		playerStatus: atomic.NewInt64(),
		volume:       atomic.NewInt64(),
		mute:         atomic.NewBool(),
	}

	device.SetSerialNumber(sn)
	device.setPlayerStatus(StatusStopped)
	device.SetVolume(config.Volume)
	device.SetMute(config.Mute)

	return device, nil
}
