package rpi

const (
	CurrentlyUnderVoltage             uint64 = 0x00001
	CurrentlyThrottled                uint64 = 0x00002
	CurrentlyARMFrequencyCapped       uint64 = 0x00004
	CurrentlySoftTemperatureReached   uint64 = 0x00008
	SinceRebootUnderVoltage           uint64 = 0x10000
	SinceRebootThrottled              uint64 = 0x20000
	SinceRebootARMFrequencyCapped     uint64 = 0x40000
	SinceRebootSoftTemperatureReached uint64 = 0x80000
)

type Throttled uint64

func (t Throttled) Is(v uint64) bool {
	return uint64(t)&v == 1
}

func (t Throttled) IsCurrentlyUnderVoltage() bool {
	return t.Is(CurrentlyUnderVoltage)
}

func (t Throttled) IsCurrentlyThrottled() bool {
	return t.Is(CurrentlyThrottled)
}

func (t Throttled) IsCurrentlyARMFrequencyCapped() bool {
	return t.Is(CurrentlyARMFrequencyCapped)
}

func (t Throttled) IsCurrentlySoftTemperatureReached() bool {
	return t.Is(CurrentlySoftTemperatureReached)
}

func (t Throttled) IsSinceRebootUnderVoltage() bool {
	return t.Is(SinceRebootUnderVoltage)
}

func (t Throttled) IsSinceRebootThrottled() bool {
	return t.Is(SinceRebootThrottled)
}

func (t Throttled) IsSinceRebootARMFrequencyCapped() bool {
	return t.Is(SinceRebootARMFrequencyCapped)
}

func (t Throttled) IsSinceRebootSoftTemperatureReached() bool {
	return t.Is(SinceRebootSoftTemperatureReached)
}
