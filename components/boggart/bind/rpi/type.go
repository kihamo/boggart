package rpi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/rpi"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	sysFS := rpi.NewSysFS()

	sn, err := sysFS.SerialNumber()
	if err != nil {
		return nil, err
	}

	config := c.(*Config)

	config.TopicModel = config.TopicModel.Format(sn)
	config.TopicCPUFrequentie = config.TopicCPUFrequentie.Format(sn)
	config.TopicTemperature = config.TopicTemperature.Format(sn)
	config.TopicVoltage = config.TopicVoltage.Format(sn)
	config.TopicCurrentlyUnderVoltage = config.TopicCurrentlyUnderVoltage.Format(sn)
	config.TopicCurrentlyThrottled = config.TopicCurrentlyThrottled.Format(sn)
	config.TopicCurrentlyARMFrequencyCapped = config.TopicCurrentlyARMFrequencyCapped.Format(sn)
	config.TopicCurrentlySoftTemperatureReached = config.TopicCurrentlySoftTemperatureReached.Format(sn)
	config.TopicSinceRebootUnderVoltage = config.TopicSinceRebootUnderVoltage.Format(sn)
	config.TopicSinceRebootThrottled = config.TopicSinceRebootThrottled.Format(sn)
	config.TopicSinceRebootARMFrequencyCapped = config.TopicSinceRebootARMFrequencyCapped.Format(sn)
	config.TopicSinceRebootSoftTemperatureReached = config.TopicSinceRebootSoftTemperatureReached.Format(sn)

	bind := &Bind{
		config:           config,
		serialNumber:     sn,
		providerVCGenCMD: rpi.NewVCGenCMD(),
		providerSysFS:    sysFS,
	}

	return bind, nil
}
