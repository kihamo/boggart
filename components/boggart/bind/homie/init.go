package homie

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/homie/esp"
)

func init() {
	boggart.RegisterBindType("homie:esp8266", esp.Type{})
	boggart.RegisterBindType("homie:esp32", esp.Type{})
}
