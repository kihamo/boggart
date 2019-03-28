package chromecast

import (
	"strconv"

	"github.com/barnybug/go-cast/log"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	log.Debug = config.Debug

	bind := &Bind{
		host: config.Host.IP,
		port: config.Port,

		volume:         atomic.NewUint32Null(),
		mute:           atomic.NewBoolNull(),
		status:         atomic.NewString(),
		mediaContentID: atomic.NewString(),

		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}
	bind.SetSerialNumber(config.Host.String() + ":" + strconv.Itoa(config.Port))

	return bind, nil
}
