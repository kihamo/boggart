package alsa

import (
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/shadow/components/dashboard"
)

type Type struct {
	dashboard.Handler
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

	bind := &Bind{
		done:         make(chan struct{}, 1),
		playerStatus: atomic.NewInt64Default(StatusStopped.Int64()),
		volume:       atomic.NewInt64Default(config.Volume),
		mute:         atomic.NewBoolDefault(config.Mute),
	}

	bind.SetSerialNumber(sn)

	return bind, nil
}
