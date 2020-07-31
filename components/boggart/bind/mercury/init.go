package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v1"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v3"
	mercury1 "github.com/kihamo/boggart/providers/mercury/v1"
)

func init() {
	boggart.RegisterBindType("mercury:200", v1.Type{
		SerialNumberFunc: mercury1.WithAddress200AsString,
		Device:           mercury1.Device200,
	})
	boggart.RegisterBindType("mercury:201", v1.Type{
		SerialNumberFunc: mercury1.WithAddressAsString,
		Device:           mercury1.Device201,
	})
	boggart.RegisterBindType("mercury:203", v1.Type{
		SerialNumberFunc: mercury1.WithAddressAsString,
		Device:           mercury1.Device203,
	})
	boggart.RegisterBindType("mercury:206", v1.Type{
		SerialNumberFunc: mercury1.WithAddressAsString,
		Device:           mercury1.Device206,
	})

	version3 := v3.Type{}
	boggart.RegisterBindType("mercury:203.2TD", version3)
	boggart.RegisterBindType("mercury:204", version3)
	boggart.RegisterBindType("mercury:208", version3)
	boggart.RegisterBindType("mercury:230", version3)
	boggart.RegisterBindType("mercury:231", version3)
	boggart.RegisterBindType("mercury:233", version3)
	boggart.RegisterBindType("mercury:234", version3)
	boggart.RegisterBindType("mercury:236", version3)
	boggart.RegisterBindType("mercury:238", version3)
}
