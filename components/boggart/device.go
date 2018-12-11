package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/event"
	"github.com/kihamo/snitch"
)

var (
	DeviceEventSyslogReceive                   = event.NewBaseEvent("SyslogReceive")
	DeviceEventDevicesManagerReady             = event.NewBaseEvent("DevicesManagerReady")
	DeviceEventDeviceRegister                  = event.NewBaseEvent("DeviceRegister")
	DeviceEventDeviceDisabledAfterCheck        = event.NewBaseEvent("DeviceDisabledAfterCheck")
	DeviceEventDeviceEnabledAfterCheck         = event.NewBaseEvent("DeviceEnabledAfterCheck")
	DeviceEventDeviceEnabled                   = event.NewBaseEvent("DeviceEnabled")
	DeviceEventDeviceDisabled                  = event.NewBaseEvent("DeviceDisabled")
	DeviceEventWifiClientConnected             = event.NewBaseEvent("WifiClientConnected")
	DeviceEventWifiClientDisconnected          = event.NewBaseEvent("WifiClientDisconnected")
	DeviceEventVPNClientConnected              = event.NewBaseEvent("VPNClientConnected")
	DeviceEventVPNClientDisconnected           = event.NewBaseEvent("VPNClientDisconnected")
	DeviceEventHikvisionEventNotificationAlert = event.NewBaseEvent("HikvisionEventNotificationAlert")
	DeviceEventSoftVideoBalanceChanged         = event.NewBaseEvent("SoftVideoBalanceChanged")
	DeviceEventMegafonBalanceChanged           = event.NewBaseEvent("SoftMegafonBalanceChanged")
	DeviceEventPulsarChanged                   = event.NewBaseEvent("PulsarChanged")
	DeviceEventPulsarPulsedChanged             = event.NewBaseEvent("PulsarPulsedChanged")
	DeviceEventMercury200Changed               = event.NewBaseEvent("Mercury200Changed")
	DeviceEventBME280Changed                   = event.NewBaseEvent("BME280")
	DeviceEventGPIOPinChanged                  = event.NewBaseEvent("GPIOPinChanged")
	DeviceEventDS18B20Changed                  = event.NewBaseEvent("DS18B20Changed")
	DeviceEventSocketStateChanged              = event.NewBaseEvent("SocketStateChanged")
	DeviceEventSocketPowerChanged              = event.NewBaseEvent("SocketPowerChanged")
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
	DeviceTypeRemoteControll
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
	IsEnabled() bool
	Disable() error
	Enable() error
	Ping(context.Context) bool
	TriggerEventChannel() <-chan DeviceTriggerEvent
}

type DeviceHasTasks interface {
	Tasks() []workers.Task
}

type DeviceHasListeners interface {
	Listeners() []workers.ListenerWithEvents
}

type DeviceHasMQTTSubscribers mqtt.HasSubscribers
