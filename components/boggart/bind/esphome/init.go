package esphome

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/mqtt"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/native_api"
)

func init() {
	nativeAPI := nativeapi.Type{}
	boggart.RegisterBindType("esphome", nativeAPI)
	boggart.RegisterBindType("esphome:native_api", nativeAPI)

	boggart.RegisterBindType("esphome:mqtt", mqtt.Type{})
}
