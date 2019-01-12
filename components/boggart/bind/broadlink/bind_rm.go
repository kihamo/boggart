package broadlink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

const (
	RMCaptureDuration = time.Second * 15
)

type BindRM struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider *broadlink.RMProPlus
}
