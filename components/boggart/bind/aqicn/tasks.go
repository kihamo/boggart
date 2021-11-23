package aqicn

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/aqicn/client/geo_localized"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	cfg := b.config()

	params := geo_localized.NewGetGeoLocalizedFeedParamsWithContext(ctx).
		WithLat(cfg.Latitude).
		WithLng(cfg.Longitude)

	response, err := b.client.GeoLocalized.GetGeoLocalizedFeed(params, nil)
	if err != nil {
		return err
	}

	id := b.Meta().ID()
	current := response.GetPayload().Data.Iaqi

	if current.T != nil {
		metricCurrentTemperature.With("id", id).Set(current.T.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentTemperature.Format(id), current.T.V); e != nil {
			err = multierr.Append(err, fmt.Errorf("send mqtt message about current temperature failed: %w", e))
		}
	}

	if current.P != nil {
		metricCurrentPressure.With("id", id).Set(current.P.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentPressure.Format(id), current.P.V); e != nil {
			err = multierr.Append(err, fmt.Errorf("send mqtt message about current pressure failed: %w", e))
		}
	}

	if current.H != nil {
		metricCurrentHumidity.With("id", id).Set(current.H.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentHumidity.Format(id), current.H.V); e != nil {
			err = multierr.Append(err, fmt.Errorf("send mqtt message about current humidity failed: %w", e))
		}
	}

	if current.Dew != nil {
		metricCurrentDewPoint.With("id", id).Set(current.Dew.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentDewPoint.Format(id), current.Dew.V); e != nil {
			err = multierr.Append(err, fmt.Errorf("send mqtt message about current dew point failed: %w", e))
		}
	}

	if current.W != nil {
		metricCurrentWindSpeed.With("id", id).Set(current.W.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentWindSpeed.Format(id), current.W.V); e != nil {
			err = multierr.Append(err, fmt.Errorf("send mqtt message about current wind speed failed: %w", e))
		}
	}

	return err
}
