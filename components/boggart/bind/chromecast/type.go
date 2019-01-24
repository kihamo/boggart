package chromecast

import (
	"strconv"

	"github.com/barnybug/go-cast/log"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	log.Debug = config.Debug

	device := &Bind{
		host: config.Host.IP,
		port: config.Port,

		volume:         atomic.NewUint32Null(),
		mute:           atomic.NewBoolNull(),
		status:         atomic.NewString(),
		mediaContentID: atomic.NewString(),

		livenessInterval: config.LivenessInterval,
		livenessTimeout:  config.LivenessTimeout,
	}
	device.SetSerialNumber(config.Host.String() + ":" + strconv.Itoa(config.Port))

	if err := device.initCast(); err != nil {
		return nil, err
	}

	return device, nil
}
