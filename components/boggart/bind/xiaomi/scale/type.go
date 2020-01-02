package scale

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/bluetooth"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	device, err := bluetooth.NewDevice()
	if err != nil {
		return nil, err
	}

	provider, err := scale.New(device, config.MAC.HardwareAddr, config.CaptureDuration)
	if err != nil {
		return nil, err
	}

	sn := config.MAC.String()

	config.TopicWeight = config.TopicWeight.Format(sn)

	bind := &Bind{
		config:   config,
		provider: provider,
	}
	bind.SetSerialNumber(sn)

	return bind, nil
}
