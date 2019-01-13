package boggart

import (
	"errors"
	"sync"
)

var (
	bindTypesMutex sync.RWMutex
	bindTypes      = make(map[string]BindType)
)

func RegisterBindType(name string, kind BindType) {
	bindTypesMutex.Lock()
	defer bindTypesMutex.Unlock()

	if kind == nil {
		panic("Bind type name is nil")
	}

	if _, dup := bindTypes[name]; dup {
		panic("Register called twice for bind type " + name)
	}

	bindTypes[name] = kind
}

func GetBindType(name string) (BindType, error) {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	kind, ok := bindTypes[name]
	if !ok {
		return nil, errors.New("Bind type " + name + " isn't register")
	}

	return kind, nil
}

func GetBindTypes() map[string]BindType {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	return bindTypes
}

type BindType interface {
	Config() interface{}
	CreateBind(config interface{}) (Bind, error)
}
