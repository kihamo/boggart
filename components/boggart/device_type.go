package boggart

import (
	"errors"
	"sync"
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

func GetDeviceTypes() map[string]DeviceType {
	deviceTypesMutex.RLock()
	defer deviceTypesMutex.RUnlock()

	return deviceTypes
}

type DeviceType interface {
	Config() interface{}
	CreateBind(config interface{}) (DeviceBind, error)
}
