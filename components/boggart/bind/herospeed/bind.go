package herospeed

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/herospeed"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
	client *herospeed.Client
}
