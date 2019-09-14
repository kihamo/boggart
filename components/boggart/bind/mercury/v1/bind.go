package v1

import (
	"github.com/kihamo/boggart/components/boggart"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config   *Config
	provider *mercury.MercuryV1
}
