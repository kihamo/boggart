package ledwifi

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("led_wifi", Type{})
}
