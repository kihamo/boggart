package openweathermap

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/openweathermap/client/weather"
	"go.uber.org/multierr"
)

const (
	TaskNameLocationCheck = "location-check"
	TaskNameUpdater       = "updater"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName(TaskNameLocationCheck).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskLocationCheck),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskLocationCheck(ctx context.Context) (err error) {
	cfg := b.config()

	b.locationMutex.Lock()
	defer b.locationMutex.Unlock()

	switch {
	case cfg.CityID > 0:
		var response *weather.GetCurrentByCityIDOK

		params := weather.NewGetCurrentByCityIDParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithID(cfg.CityID)

		response, err = b.client.Weather.GetCurrentByCityID(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}

	case cfg.CityName != "":
		var response *weather.GetCurrentByCityNameOK

		params := weather.NewGetCurrentByCityNameParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithQ(cfg.CityName)

		response, err = b.client.Weather.GetCurrentByCityName(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}

	case cfg.Latitude != 0 && cfg.Longitude != 0:
		var response *weather.GetCurrentByGeographicCoordinatesOK

		params := weather.NewGetCurrentByGeographicCoordinatesParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithLat(cfg.Latitude).
			WithLon(cfg.Longitude)

		response, err = b.client.Weather.GetCurrentByGeographicCoordinates(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}

	case cfg.Zip != "":
		var response *weather.GetCurrentByZIPCodeOK

		params := weather.NewGetCurrentByZIPCodeParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithZip(cfg.Zip)

		response, err = b.client.Weather.GetCurrentByZIPCode(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}

	default:
		err = errors.New("location is empty")
	}

	if err != nil {
		return err
	}

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)
	if err != nil {
		return fmt.Errorf("register task "+TaskNameUpdater+" failed: %w", err)
	}

	return nil
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	response, err := b.OneCallMigrate(ctx)
	if err != nil {
		return err
	}

	id := b.Meta().ID()
	cfg := b.config()

	metricCurrent.With("id", id).Set(response.Current.Temp)
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicCurrentTemp.Format(id), response.Current.Temp); e != nil {
		err = multierr.Append(err, e)
	}

	for i, day := range response.Daily {
		dayAsString := strconv.Itoa(i)

		metricTempMin.With("id", id).With("day", dayAsString).Set(day.Temp.Min)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyTempMin.Format(id, i), day.Temp.Min); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempMax.With("id", id).With("day", dayAsString).Set(day.Temp.Max)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyTempMax.Format(id, i), day.Temp.Max); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempDay.With("id", id).With("day", dayAsString).Set(day.Temp.Day)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyTempDay.Format(id, i), day.Temp.Day); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempNight.With("id", id).With("day", dayAsString).Set(day.Temp.Night)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyTempNight.Format(id, i), day.Temp.Night); e != nil {
			err = multierr.Append(err, e)
		}

		metricTempMorning.With("id", id).With("day", dayAsString).Set(day.Temp.Morn)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyTempMorning.Format(id, i), day.Temp.Morn); e != nil {
			err = multierr.Append(err, e)
		}

		metricWindSpeed.With("id", id).With("day", dayAsString).Set(day.WindSpeed)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDailyWindSpeed.Format(id, i), day.WindSpeed); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
