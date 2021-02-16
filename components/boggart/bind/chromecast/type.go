package chromecast

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		disconnected:   atomic.NewBoolNull(),
		volume:         atomic.NewUint32Null(),
		mute:           atomic.NewBoolNull(),
		status:         atomic.NewString(),
		mediaContentID: atomic.NewString(),
	}
}
