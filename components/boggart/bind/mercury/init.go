package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v1"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v3"
)

func init() {
	boggart.RegisterBindType("mercury:200", v1.Type{})
	boggart.RegisterBindType("mercury:230", v3.Type{})
}
