package aqicn

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/aqicn"
	"github.com/kihamo/boggart/providers/aqicn/client/feed"
	"github.com/kihamo/boggart/providers/aqicn/models"
)

// TODO: посчитать в скольки километрах от места находится станция наблюдения

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.WidgetBind
	di.WorkersBind

	client *aqicn.Client
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() (err error) {
	cfg := b.config()

	b.client = aqicn.New(cfg.Token, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	return nil
}

func (b *Bind) Feed(ctx context.Context) (f *models.FeedData, err error) {
	var response interface {
		GetPayload() *models.Feed
	}

	cfg := b.config()

	switch {
	case cfg.CityName != "":
		params := feed.NewGetByCityParamsWithContext(ctx).
			WithCity(cfg.CityName)

		response, err = b.client.Feed.GetByCity(params, nil)
	case cfg.Latitude != 0 && cfg.Longitude != 0:
		params := feed.NewGetByLatLngParamsWithContext(ctx).
			WithLat(cfg.Latitude).
			WithLng(cfg.Longitude)

		response, err = b.client.Feed.GetByLatLng(params, nil)
	default:
		err = errors.New("location is empty")
	}

	if err == nil {
		f = response.GetPayload().Data
	}

	return f, err
}
