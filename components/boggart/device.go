package boggart

import (
	"context"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/event"
	"github.com/kihamo/snitch"
)

var (
	DeviceEventSyslogReceive            = event.NewBaseEvent("SyslogReceive")
	DeviceEventDevicesManagerReady      = event.NewBaseEvent("DevicesManagerReady")
	DeviceEventDeviceRegister           = event.NewBaseEvent("DeviceRegister")
	DeviceEventDeviceDisabledAfterCheck = event.NewBaseEvent("DeviceDisabledAfterCheck")
	DeviceEventDeviceEnabledAfterCheck  = event.NewBaseEvent("DeviceEnabledAfterCheck")
	DeviceEventWifiClientConnected      = event.NewBaseEvent("WifiClientConnected")
	DeviceEventWifiClientDisconnected   = event.NewBaseEvent("WifiClientDisconnected")
)

type DeviceId int64

const (
	DeviceIdElectricityMeter DeviceId = iota
	DeviceIdCameraHall
	DeviceIdCameraStreet
	DeviceIdHeatMeter
	DeviceIdEntranceDoor
	DeviceIdPhone
	DeviceIdRouter
	DeviceIdUPS
	DeviceIdVideoRecorder
	DeviceIdWaterMeterCold
	DeviceIdWaterMeterHot
)

type DeviceType int64

const (
	DeviceTypeCamera DeviceType = iota
	DeviceTypeDoor
	DeviceTypeElectricityMeter
	DeviceTypeHeatMeter
	DeviceTypeInternetProvider
	DeviceTypePhone
	DeviceTypeRouter
	DeviceTypeUPS
	DeviceTypeVideoRecorder
	DeviceTypeWaterMeter
)

type DevicesManager interface {
	snitch.Collector

	Register(Device) string
	RegisterWithID(string, Device)
	Device(string) Device
	Devices() map[string]Device
	DevicesByTypes([]DeviceType) map[string]Device
	Attach(workers.Event, workers.Listener)
	DeAttach(workers.Event, workers.Listener)
	Listeners() []workers.Listener
	GetListenerMetadata(id string) workers.Metadata
}

type DeviceTriggerEvent interface {
	Event() workers.Event
	Arguments() []interface{}
}

type Device interface {
	Id() string
	Description() string
	Types() []DeviceType
	IsEnabled() bool
	Disable()
	Enable()
	Ping(context.Context) bool
	Tasks() []workers.Task
	Listeners() []workers.ListenerWithEvents
	TriggerEventChannel() <-chan DeviceTriggerEvent
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

	Number() string
	Balance(context.Context) (float64, error)
}

type ElectricityMeter interface {
	Device

	Tariffs(context.Context) (map[string]float64, error)
	Voltage(context.Context) (float64, error)
	Amperage(context.Context) (float64, error)
	Power(context.Context) (float64, error)
}

type WaterMeter interface {
	Device

	Volume(context.Context) (float64, error)
}

type HeatMeter interface {
	Device

	TemperatureIn(context.Context) (float64, error)
	TemperatureOut(context.Context) (float64, error)
	TemperatureDelta(context.Context) (float64, error)
	Energy(context.Context) (float64, error)
	Consumption(context.Context) (float64, error)
}

type Router interface {
	Device

	// TODO: public interface
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
