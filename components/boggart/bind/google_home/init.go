package google_home

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("google_home:mini", Type{})
}
