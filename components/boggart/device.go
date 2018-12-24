package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
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
	DeviceEventDeviceEnabled            = event.NewBaseEvent("DeviceEnabled")
	DeviceEventDeviceDisabled           = event.NewBaseEvent("DeviceDisabled")
)

type DeviceId int64

const (
	DeviceIdElectricityMeter DeviceId = iota
	DeviceIdHeatMeter
	DeviceIdPhone
	DeviceIdWaterMeterCold
	DeviceIdWaterMeterHot
)

type DeviceType int64

const (
	DeviceTypeElectricityMeter DeviceType = iota
	DeviceTypeHeatMeter
	DeviceTypeInternetProvider
	DeviceTypePhone
	DeviceTypeRouter
	DeviceTypeCamera
	DeviceTypeWaterMeter
	DeviceTypeThermometer
	DeviceTypeBarometer
	DeviceTypeHygrometer
	DeviceTypeGPIO
	DeviceTypeSocket
	DeviceTypeRemoteControl
	DeviceTypeLED
	DeviceTypeTV
)

type DeviceStatus uint64

const (
	DeviceStatusUnknown DeviceStatus = iota
	DeviceStatusUninitialized
	DeviceStatusInitializing
	DeviceStatusOnline
	DeviceStatusOffline
	DeviceStatusRemoving
	DeviceStatusRemoved
)

type DevicesManager interface {
	snitch.Collector

	Register(Device) string
	RegisterWithID(string, Device)
	Device(string) Device
	Devices() map[string]Device
	DevicesByTypes([]DeviceType) map[string]Device
	Check()
	CheckByKeys(...string)
	IsReady() bool
}

type DeviceTriggerEvent interface {
	Context() context.Context
	Event() workers.Event
	Arguments() []interface{}
}

type Device interface {
	Id() string
	Description() string
	Types() []DeviceType
	Status() DeviceStatus

	// deprecated
	IsEnabled() bool
	// deprecated
	Disable() error
	// deprecated
	Enable() error
	// deprecated
	Ping(context.Context) bool
	// deprecated
	TriggerEventChannel() <-chan DeviceTriggerEvent
}

type DeviceHasSerialNumber interface {
	SerialNumber() string
}

type DeviceHasTasks interface {
	Tasks() []workers.Task
}

type DeviceHasListeners interface {
	Listeners() []workers.ListenerWithEvents
}

type DeviceHasMQTTClient interface {
	SetMQTTClient(mqtt.Component)
}

type DeviceHasMQTTSubscribers mqtt.HasSubscribers

type DeviceHasMQTTTopics interface {
	MQTTTopics() []mqtt.Topic
}
