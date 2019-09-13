package elektroset

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
	client *elektroset.Client
}
