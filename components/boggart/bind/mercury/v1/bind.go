package v1

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	provider        *mercury.MercuryV1
	updaterInterval time.Duration
}
