package rkcm

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/rkcm"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
	client *rkcm.Client
}
