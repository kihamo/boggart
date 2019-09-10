package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/rkcm"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *rkcm.Client
	config *Config
}
