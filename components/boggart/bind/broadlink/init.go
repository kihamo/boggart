package broadlink

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("broadlink_rm", TypeRM{})
	boggart.RegisterDeviceType("broadlink_sp3s", TypeSP3S{})
}
