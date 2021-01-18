package homie

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/homie/esp"
)

func init() {
	boggart.RegisterBindType("homie", esp.Type{}, "homie:esp8266", "homie:esp32")
}
