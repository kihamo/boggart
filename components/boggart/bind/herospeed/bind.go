package herospeed

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/herospeed"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *herospeed.Client
	config *Config
}
