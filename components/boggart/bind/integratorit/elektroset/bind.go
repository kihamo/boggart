package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/integratorit/elektroset"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *elektroset.Client
	config *Config
}
