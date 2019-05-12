package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
	device *devices.Vacuum

	battery   *atomic.Uint32Null
	cleanArea *atomic.Uint32Null
	cleanTime *atomic.Uint32Null
}

func (b *Bind) Close() error {
	return b.device.Close()
}
