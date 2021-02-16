package sp3s

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/broadlink"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	state *atomic.BoolNull
	power *atomic.Float32Null

	provider *broadlink.SP3S
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	switch cfg.Model {
	case "sp3seu":
		b.provider = broadlink.NewSP3SEU(cfg.MAC.HardwareAddr, cfg.Host)

	case "sp3sus":
		b.provider = broadlink.NewSP3SUS(cfg.MAC.HardwareAddr, cfg.Host)

	default:
		return errors.New("unknown model " + cfg.Model)
	}

	b.provider.SetTimeout(cfg.ConnectionTimeout)

	b.Meta().SetMAC(cfg.MAC.HardwareAddr)
	b.state.Nil()
	b.power.Nil()

	return nil
}

func (b *Bind) State() (bool, error) {
	return b.provider.State()
}

func (b *Bind) On(ctx context.Context) error {
	err := b.provider.On()
	if err == nil {
		err = b.taskUpdaterHandler(ctx)
	}

	return err
}

func (b *Bind) Off(ctx context.Context) error {
	err := b.provider.Off()
	if err == nil {
		err = b.taskUpdaterHandler(ctx)
	}

	return err
}

func (b *Bind) Power() (float64, error) {
	return b.provider.Power()
}
