package influxdb

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	config *Config
	client influxdb2.Client
}

func (b *Bind) Close() error {
	b.client.Close()
	return nil
}
