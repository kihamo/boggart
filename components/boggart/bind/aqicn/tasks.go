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

	const mqttErrTpl = "send mqtt message about %s failed: %w"

	if current.T != nil {
		metricCurrentTemperature.With("id", id).Set(current.T.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentTemperature.Format(id), current.T.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current temperatur", e))
		}
	}

	if current.P != nil {
		metricCurrentPressure.With("id", id).Set(current.P.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentPressure.Format(id), current.P.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current pressure", e))
		}
	}

	if current.H != nil {
		metricCurrentHumidity.With("id", id).Set(current.H.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentHumidity.Format(id), current.H.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current humidity", e))
		}
	}

	if current.Dew != nil {
		metricCurrentDewPoint.With("id", id).Set(current.Dew.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentDewPoint.Format(id), current.Dew.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current dew point", e))
		}
	}

	if current.W != nil {
		metricCurrentWindSpeed.With("id", id).Set(current.W.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentWindSpeed.Format(id), current.W.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current wind speed", e))
		}
	}

	if current.Pm25 != nil {
		metricCurrentPm25Value.With("id", id).Set(current.Pm25.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentPm25Value.Format(id), current.Pm25.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current PM25 value", e))
		}
	}

	if current.Pm10 != nil {
		metricCurrentPm10Value.With("id", id).Set(current.Pm10.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentPm10Value.Format(id), current.Pm10.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current PM10 value", e))
		}
	}

	if current.O3 != nil {
		metricCurrentO3Value.With("id", id).Set(current.O3.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentO3Value.Format(id), current.O3.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current O3 value", e))
		}
	}

	if current.No2 != nil {
		metricCurrentNO2Value.With("id", id).Set(current.No2.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentNO2Value.Format(id), current.No2.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current NO2 value", e))
		}
	}

	if current.Co != nil {
		metricCurrentCOValue.With("id", id).Set(current.Co.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentCOValue.Format(id), current.Co.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current CO value", e))
		}
	}

	if current.So2 != nil {
		metricCurrentSO2Value.With("id", id).Set(current.So2.V)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentSO2Value.Format(id), current.So2.V); e != nil {
			err = multierr.Append(err, fmt.Errorf(mqttErrTpl, "current SO2 value", e))
		}
	}

	return err
}
