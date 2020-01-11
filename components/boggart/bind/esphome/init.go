package esphome

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/mqtt"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/native_api"
)

func init() {
	nativeApi := native_api.Type{}
	boggart.RegisterBindType("esphome", nativeApi)
	boggart.RegisterBindType("esphome:native_api", nativeApi)

	boggart.RegisterBindType("esphome:mqtt", mqtt.Type{})
}
