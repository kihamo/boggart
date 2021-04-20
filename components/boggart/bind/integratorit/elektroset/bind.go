package elektroset

import (
	"context"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
)

const (
	layoutPeriod = "2006-01-02"
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

	client      *elektroset.Client
	metersCount *atomic.Uint32Null
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.client = elektroset.New(cfg.Login, cfg.Password)
	b.metersCount.Nil()

	return nil
}

func (b *Bind) Houses(ctx context.Context) ([]elektroset.House, error) {
	houses, err := b.client.Houses(ctx)
	if houseID := b.config().HouseID; err == nil && houseID > 0 {
		for _, house := range houses {
			if house.ID == houseID {
				return []elektroset.House{house}, nil
			}
		}

		return nil, nil
	}

	return houses, err
}
