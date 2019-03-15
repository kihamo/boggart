package sun

import (
	"context"
	"time"

	"github.com/kelvins/sunrisesunset"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	b.taskStateUpdater = task.NewFunctionTask(b.taskUpdater)
	b.taskStateUpdater.SetRepeats(-1)
	b.taskStateUpdater.SetName("bind-astro-sun-updater-" + b.SerialNumber())

	return []workers.Task{
		b.taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(b.config.Lat, b.config.Lon, b.config.UTCOffset, today)
	if err != nil {
		return nil, err
	}

	serialNumberMQTT := mqtt.NameReplace(b.SerialNumber())

	sunrise = time.Date(today.Year(), today.Month(), today.Day(), sunrise.Hour(), sunrise.Minute(), sunrise.Second(), sunrise.Nanosecond(), today.Location())
	sunset = time.Date(today.Year(), today.Month(), today.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), today.Location())

	if ok := b.sunrise.Set(sunrise); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSunrise.Format(serialNumberMQTT), sunrise); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.sunset.Set(sunset); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSunset.Format(serialNumberMQTT), sunset); e != nil {
			err = multierr.Append(err, e)
		}
	}

	dayLight := sunset.Sub(sunrise)
	if ok := b.dayLight.Set(dayLight); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicDayLight.Format(serialNumberMQTT), dayLight); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// change start date (only one of 24 hours)
	b.taskStateUpdater.SetRepeatInterval(today.Add(time.Hour * 24).Sub(now))

	return nil, err
}
