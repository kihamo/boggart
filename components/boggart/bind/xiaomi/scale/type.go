package scale

import (
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	mac := config.MAC.String()

	config.TopicProfile = config.TopicProfile.Format(mac)
	config.TopicProfileActivate = config.TopicProfileActivate.Format(mac)
	config.TopicProfileSettings = config.TopicProfileSettings.Format(mac)
	config.TopicDatetime = config.TopicDatetime.Format(mac)
	config.TopicWeight = config.TopicWeight.Format(mac)
	config.TopicImpedance = config.TopicImpedance.Format(mac)
	config.TopicBMR = config.TopicBMR.Format(mac)
	config.TopicBMI = config.TopicBMI.Format(mac)
	config.TopicFatPercentage = config.TopicFatPercentage.Format(mac)
	config.TopicWaterPercentage = config.TopicWaterPercentage.Format(mac)
	config.TopicIdealWeight = config.TopicIdealWeight.Format(mac)
	config.TopicLBMCoefficient = config.TopicLBMCoefficient.Format(mac)
	config.TopicBoneMass = config.TopicBoneMass.Format(mac)
	config.TopicMuscleMass = config.TopicMuscleMass.Format(mac)
	config.TopicVisceralFat = config.TopicVisceralFat.Format(mac)
	config.TopicFatMassToIdeal = config.TopicFatMassToIdeal.Format(mac)
	config.TopicProteinPercentage = config.TopicProteinPercentage.Format(mac)
	config.TopicBodyType = config.TopicBodyType.Format(mac)
	config.TopicMetabolicAge = config.TopicMetabolicAge.Format(mac)

	bind := &Bind{
		disconnected:         atomic.NewBoolNull(),
		config:               config,
		measureStartDatetime: atomic.NewTimeDefault(time.Now()),
	}

	if len(config.Profiles) > 0 {
		for name, profile := range config.Profiles {
			profile.Name = name
			profile.Age = profile.GetAge()
		}
	}

	return bind, nil
}
