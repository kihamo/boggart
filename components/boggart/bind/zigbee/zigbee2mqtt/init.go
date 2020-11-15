package zigbee2mqtt

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("zigbee:zigbee2mqtt", Type{})
}
