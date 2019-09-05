package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mosenergosbyt"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *mosenergosbyt.Client
	config *Config
}
