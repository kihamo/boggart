package google_home_mini

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("google_home_mini", Type{})
}
