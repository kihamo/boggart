package boggart

import (
	"context"
	"fmt"
	"strings"

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

type DeviceHasMQTTClient interface {
	SetMQTTClient(mqtt.Component)
}

type DeviceHasMQTTSubscribers mqtt.HasSubscribers

type DeviceMQTTTopic string

type DeviceHasMQTTTopics interface {
	MQTTTopics() []DeviceMQTTTopic
}

func (t DeviceMQTTTopic) String() string {
	return string(t)
}

func (t DeviceMQTTTopic) Format(args ...interface{}) string {
	parts := mqtt.RouteSplit(t.String())

	for _, arg := range args {
		for i, topic := range parts {
			if topic == "+" {
				parts[i] = fmt.Sprintf("%v", arg)
				break
			}
		}
	}

	return strings.Join(parts, "/")
}
