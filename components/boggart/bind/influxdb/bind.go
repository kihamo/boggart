package influxdb

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	client influxdb2.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	authToken := cfg.AuthToken
	if authToken == "" && cfg.DSN.User != nil {
		authToken = cfg.DSN.User.String()
	}

	scheme := cfg.DSN.Scheme
	if scheme == "" {
		scheme = "http"
	}

	b.client = influxdb2.NewClient(scheme+"://"+cfg.DSN.Host, authToken)

	return nil
}

func (b *Bind) Close() error {
	b.client.Close()
	return nil
}
