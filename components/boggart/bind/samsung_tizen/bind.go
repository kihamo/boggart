package tizen

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	config *Config
	client *tv.APIv2
}

func (b *Bind) Run() error {
	if b.config.MAC != nil && b.Meta().MAC() == nil {
		b.Meta().SetMAC(b.config.MAC.HardwareAddr)
	}

	return nil
}
