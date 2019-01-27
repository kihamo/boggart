package broadlink

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/broadlink/rm"
	"github.com/kihamo/boggart/components/boggart/bind/broadlink/sp3s"
)

func init() {
	boggart.RegisterBindType("broadlink:rm", rm.Type{})
	boggart.RegisterBindType("broadlink:sp3s", sp3s.Type{})
}
