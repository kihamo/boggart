package alsa

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		playerStatus: atomic.NewInt64(),
		volume:       atomic.NewInt64(),
		mute:         atomic.NewBool(),
	}
}
