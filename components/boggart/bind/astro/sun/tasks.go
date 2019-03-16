package sun

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/mourner/suncalc-go"
	"go.uber.org/multierr"
)

const (
	dayDuration = time.Hour * 24
)

func (b *Bind) Tasks() []workers.Task {
	b.taskStateUpdater = task.NewFunctionTask(b.taskUpdater)
	b.taskStateUpdater.SetRepeats(-1)
	b.taskStateUpdater.SetName("bind-astro-sun-updater-" + b.SerialNumber())

	return []workers.Task{
		b.taskStateUpdater,
	}
}

/*
+ nadir:2019-03-16 00:39:15 +0300 MSK
+ solarNoon:2019-03-16 12:39:15 +0300 MSK

+ night:2019-03-16 20:39:50 +0300 MSK
+ nightEnd:2019-03-16 04:38:39 +0300 MSK

goldenHourEnd:2019-03-16 07:33:51 +0300 MSK
goldenHour:2019-03-16 17:44:39 +0300 MSK

+ sunrise:2019-03-16 06:44:37 +0300 MSK
+ sunriseEnd:2019-03-16 06:48:26 +0300 MSK

+ sunsetStart:2019-03-16 18:30:04 +0300 MSK
+ sunset:2019-03-16 18:33:52 +0300 MSK

dawn:2019-03-16 06:07:42 +0300 MSK
nauticalDawn:2019-03-16 05:24:10 +0300 MSK

dusk:2019-03-16 19:10:47 +0300 MSK
nauticalDusk:2019-03-16 19:54:19 +0300 MSK
*/

func (b *Bind) taskUpdater(ctx context.Context) (_ interface{}, err error) {
	now := time.Now()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	timesToday := suncalc.SunTimes(now, b.config.Lat, b.config.Lon)
	timesTomorrow := suncalc.SunTimes(now.Add(dayDuration), b.config.Lat, b.config.Lon)

	serialNumberMQTT := mqtt.NameReplace(b.SerialNumber())

	if ok := b.riseStart.Set(timesToday["sunrise"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseStart.Format(serialNumberMQTT), timesToday["sunrise"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.riseEnd.Set(timesToday["sunriseEnd"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseEnd.Format(serialNumberMQTT), timesToday["sunriseEnd"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	riseDuration := timesToday["sunriseEnd"].Sub(timesToday["sunrise"])
	if ok := b.riseDuration.Set(riseDuration); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRiseDuration.Format(serialNumberMQTT), riseDuration); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.setStart.Set(timesToday["sunsetStart"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetStart.Format(serialNumberMQTT), timesToday["sunsetStart"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.setEnd.Set(timesToday["sunset"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetEnd.Format(serialNumberMQTT), timesToday["sunset"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	setDuration := timesToday["sunset"].Sub(timesToday["sunsetStart"])
	if ok := b.setDuration.Set(setDuration); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSetDuration.Format(serialNumberMQTT), setDuration); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.nightStart.Set(timesToday["night"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightStart.Format(serialNumberMQTT), timesToday["night"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.nightEnd.Set(timesToday["nightEnd"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightEnd.Format(serialNumberMQTT), timesToday["nightEnd"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	nightDuration := todayEnd.Sub(timesToday["night"]) + timesTomorrow["nightEnd"].Sub(todayEnd)
	if ok := b.nightDuration.Set(nightDuration); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNightDuration.Format(serialNumberMQTT), nightDuration); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.nightEnd.Set(timesToday["nadir"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicNadir.Format(serialNumberMQTT), timesToday["nadir"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.nightEnd.Set(timesToday["solarNoon"]); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSolarNoon.Format(serialNumberMQTT), timesToday["solarNoon"]); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// change start date (only one of 24 hours)
	b.taskStateUpdater.SetRepeatInterval(todayStart.Add(dayDuration).Sub(now))

	return nil, err
}
