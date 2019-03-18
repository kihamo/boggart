package homie

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config

	deviceAttributes *sync.Map
}

func (b *Bind) registerDeviceAttributes(name string, value interface{}) {
	b.deviceAttributes.Store(name, value)
}

func (b *Bind) DeviceAttributes() map[string]interface{} {
	return syncMapToMap(b.deviceAttributes)
}
