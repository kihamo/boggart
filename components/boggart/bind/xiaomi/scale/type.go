package scale

import (
	"github.com/kihamo/boggart/atomic"
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

	provider, err := scale.NewClient(device, config.MAC.HardwareAddr, config.CaptureDuration)
	if err != nil {
		return nil, err
	}

	sn := config.MAC.String()

	config.TopicDatetime = config.TopicDatetime.Format(sn)
	config.TopicWeight = config.TopicWeight.Format(sn)
	config.TopicImpedance = config.TopicImpedance.Format(sn)
	config.TopicProfile = config.TopicProfile.Format(sn)
	config.TopicBMR = config.TopicBMR.Format(sn)
	config.TopicBMI = config.TopicBMI.Format(sn)
	config.TopicFatPercentage = config.TopicFatPercentage.Format(sn)
	config.TopicWaterPercentage = config.TopicWaterPercentage.Format(sn)
	config.TopicIdealWeight = config.TopicIdealWeight.Format(sn)
	config.TopicLBMCoefficient = config.TopicLBMCoefficient.Format(sn)
	config.TopicBoneMass = config.TopicBoneMass.Format(sn)
	config.TopicMuscleMass = config.TopicMuscleMass.Format(sn)
	config.TopicVisceralFat = config.TopicVisceralFat.Format(sn)
	config.TopicFatMassToIdeal = config.TopicFatMassToIdeal.Format(sn)
	config.TopicProteinPercentage = config.TopicProteinPercentage.Format(sn)
	config.TopicBodyType = config.TopicBodyType.Format(sn)
	config.TopicMetabolicAge = config.TopicMetabolicAge.Format(sn)

	bind := &Bind{
		config:   config,
		provider: provider,
		sex:      atomic.NewBoolNull(),
		height:   atomic.NewUint32Null(),
		age:      atomic.NewUint32Null(),
	}
	bind.SetSerialNumber(sn)

	return bind, nil
}
