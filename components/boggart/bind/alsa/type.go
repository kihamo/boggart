package alsa

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/kihamo/boggart/components/boggart"
	a "github.com/kihamo/boggart/components/voice/players/alsa"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	sn, err := machineid.ID()
	if err != nil {
		return nil, err
	}

	device := &Bind{
		player:          a.New(),
		updaterInterval: config.UpdaterInterval,
		status:          -1,
		volume:          -1,
		mute:            0,
	}

	device.Init()
	device.SetSerialNumber(sn)
	device.UpdateStatus(boggart.BindStatusOnline)

	return device, nil
}
