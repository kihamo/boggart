package hilink

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
	client *hilink.Client
}
