package esphome

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/mqtt"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/native_api"
)

func init() {
	boggart.RegisterBindType("esphome", nativeapi.Type{}, "esphome:native_api")
	boggart.RegisterBindType("esphome:mqtt", mqtt.Type{})
}
