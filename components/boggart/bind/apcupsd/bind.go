package apcupsd

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/apcupsd"
	"github.com/kihamo/boggart/protocols/apcupsd/file"
	"github.com/kihamo/boggart/protocols/apcupsd/nis"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	//di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	client *apcupsd.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	var statusReader, eventsReader apcupsd.Adapter

	if cfg.StatusFile != "" {
		statusReader = file.NewStatusReader(cfg.StatusFile)
	}

	if statusReader == nil && cfg.Address.Host != "" {
		statusReader = nis.NewStatusReader(cfg.Address.Host)
	}

	if statusReader == nil {
		return errors.New("status reader isn't initialized")
	}

	if cfg.EventsFile != "" {
		eventsReader = file.NewStatusReader(cfg.EventsFile)
	}

	if eventsReader == nil && cfg.Address.Host != "" {
		eventsReader = nis.NewEventsReader(cfg.Address.Host)
	}

	if eventsReader == nil {
		return errors.New("events reader isn't initialized")
	}

	b.client = apcupsd.NewClient(statusReader, eventsReader)

	return nil
}
