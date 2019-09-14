package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	mercury "github.com/kihamo/boggart/providers/mercury/v3"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config   *Config
	provider *mercury.MercuryV3
}
