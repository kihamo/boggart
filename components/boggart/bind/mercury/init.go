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

	boggart.RegisterBindType("mercury:v3", v3.Type{}, "mercury:203.2TD", "mercury:204", "mercury:208",
		"mercury:230", "mercury:231", "mercury:233", "mercury:234", "mercury:236", "mercury:238",
	)
}
