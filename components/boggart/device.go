package boggart

import (
	"errors"
	"sync"

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

var (
	deviceKindsMutex sync.RWMutex
	deviceKinds      = make(map[string]DeviceKind)
)

func RegisterKind(name string, kind DeviceKind) {
	deviceKindsMutex.Lock()
	defer deviceKindsMutex.Unlock()

	if kind == nil {
		panic("Device kind is nil")
	}

	if _, dup := deviceKinds[name]; dup {
		panic("RegisterKind called twice for device kind " + name)
	}

	deviceKinds[name] = kind
}

func GetKind(name string) (DeviceKind, error) {
	deviceKindsMutex.RLock()
	defer deviceKindsMutex.RUnlock()

	kind, ok := deviceKinds[name]
	if !ok {
		return nil, errors.New("Device kind " + name + " isn't register")
	}

	return kind, nil
}

type DeviceKind interface {
	Create(config map[string]interface{}) (Device, error)
}

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
	DeviceTypeUPS
	DeviceTypeSmartSpeaker
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
	IsReady() bool
}

type Device interface {
	Id() string
	Description() string
	Types() []DeviceType
	Status() DeviceStatus
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
