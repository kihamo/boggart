package mosoblgaz

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/mosoblgaz"
)

var (
	link, _ = url.Parse(mosoblgaz.BaseURL)
)

func init() {
	if link != nil {
		link.Path = "/"
	}
}

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind

	provider *mosoblgaz.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.Meta().SetLink(link)
	b.provider = mosoblgaz.New().WithToken(cfg.Token)

	if cfg.Debug {
		b.provider = b.provider.WithDebug(cfg.Debug)
	}

	return nil
}
