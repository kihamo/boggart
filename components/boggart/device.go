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
	deviceTypesMutex sync.RWMutex
	deviceTypes      = make(map[string]DeviceType)
)

func RegisterDeviceType(name string, kind DeviceType) {
	deviceTypesMutex.Lock()
	defer deviceTypesMutex.Unlock()

	if kind == nil {
		panic("Device type name is nil")
	}

	if _, dup := deviceTypes[name]; dup {
		panic("Register called twice for device type " + name)
	}

	deviceTypes[name] = kind
}

func GetDeviceType(name string) (DeviceType, error) {
	deviceTypesMutex.RLock()
	defer deviceTypesMutex.RUnlock()

	kind, ok := deviceTypes[name]
	if !ok {
		return nil, errors.New("Device type " + name + " isn't register")
	}

	return kind, nil
}

type DeviceType interface {
	CreateBind(config map[string]interface{}) (DeviceBind, error)
}

type DeviceId int64

const (
	DeviceIdElectricityMeter DeviceId = iota
	DeviceIdHeatMeter
	DeviceIdWaterMeterCold
	DeviceIdWaterMeterHot
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

	Register(device DeviceBind, t string, description string, tags []string, config map[string]interface{}) string
	RegisterWithID(id string, bind DeviceBind, t string, description string, tags []string, config map[string]interface{})
	Device(id string) Device
	Devices() map[string]Device
	IsReady() bool
}

type Device interface {
	Bind() DeviceBind
	Id() string
	Type() string
	Description() string
	Tags() []string
	Config() map[string]interface{}
}

type DeviceBind interface {
	Status() DeviceStatus
}

type DeviceBindHasSerialNumber interface {
	SerialNumber() string
}

type DeviceBindHasTasks interface {
	Tasks() []workers.Task
}

type DeviceBindHasListeners interface {
	Listeners() []workers.ListenerWithEvents
}

type DeviceBindHasMQTTClient interface {
	SetMQTTClient(mqtt.Component)
}

type DeviceBindHasMQTTSubscribers mqtt.HasSubscribers

type DeviceBindHasMQTTTopics interface {
	MQTTTopics() []mqtt.Topic
}
