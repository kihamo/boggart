package integratorit

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/integratorit/elektroset"
	"github.com/kihamo/boggart/components/boggart/bind/integratorit/mosenergosbyt"
)

func init() {
	boggart.RegisterBindType("integratorit:elektroset", elektroset.Type{})
	boggart.RegisterBindType("integratorit:mosenergosbyt", mosenergosbyt.Type{})
}
