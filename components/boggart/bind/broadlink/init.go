package broadlink

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("broadlink:rm", TypeRM{})
	boggart.RegisterBindType("broadlink:sp3s", TypeSP3S{})
}
