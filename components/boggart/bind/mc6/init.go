package mc6

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/mc6/mqtt"
)

func init() {
	boggart.RegisterBindType("mc6:mqtt", mqtt.Type{})
}
