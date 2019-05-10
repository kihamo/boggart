package miio

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/devices"
)

type Bind struct {
	boggart.BindBase

	config *Config
	device *devices.Vacuum
}
