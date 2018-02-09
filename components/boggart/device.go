package boggart

import (
	"context"

	"github.com/kihamo/snitch"
)

type DeviceManager interface {
	snitch.Collector

	Register(string, Device)
	Device(string) Device
	Devices() map[string]Device
}

type Device interface {
	Id() string
	Position() (int64, int64)
	Description() string
	IsEnabled() bool
}

type ReedSwitch interface {
	Device

	IsOpen() bool
	IsClose() bool
}

type Camera interface {
	Device

	Snapshot(context.Context) ([]byte, error)
}

type Phone interface {
	Device

	Number() uint64
}

type ElectricityMeter interface {
	Device
}

type WaterMeter interface {
	Device
}

type HeatMeter interface {
	Device
}

type Router interface {
	Device
}

type Thermometer interface {
	Device
}

type Barometer interface {
	Device
}

type WaterDetector interface {
	Device
}

type MotionDetector interface {
	Device
}
