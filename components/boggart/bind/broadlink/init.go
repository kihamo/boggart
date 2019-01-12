package broadlink

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("broadlink:rm", TypeRM{})
	boggart.RegisterDeviceType("broadlink:sp3s", TypeSP3S{})
}
