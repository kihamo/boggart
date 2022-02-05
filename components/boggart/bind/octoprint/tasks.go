package octoprint

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/octoprint/client/job"
	"github.com/kihamo/boggart/providers/octoprint/client/plugin"
	"github.com/kihamo/boggart/providers/octoprint/client/printer"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("settings").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSettingsHandler),
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

func (b *Bind) taskSettingsHandler(ctx context.Context) (err error) {
	if err = b.SystemSettingsUpdate(ctx); err != nil {
		return err
	}

	cfg := b.config()
	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
						cfg.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)
	if err != nil {
		return err
	}

	subscribers := make([]mqtt.Subscriber, 0, 2)

	if cfg := b.PluginMQTTSettings(); cfg != nil && cfg.Publish.EventActive && cfg.Publish.Events.Settings && cfg.Publish.EventTopic != "" {
		topic := mqtt.Topic(cfg.Publish.BaseTopic +
			strings.Replace(cfg.Publish.EventTopic, "{event}", "SettingsUpdated", 1))

		subscribers = append(subscribers, mqtt.NewSubscriber(topic, 0,
			b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.SystemSettingsUpdate(ctx)
			})))
	}

	// temperature
	if b.TemperatureFromMQTT() {
		topic := b.TemperatureTopic()

		parts := topic.Split()
		offset := len(parts) - 1
		for ; offset >= 0; offset-- {
			if parts[offset] == "+" {
				break
			}
		}

		subscribers = append(subscribers, mqtt.NewSubscriber(topic, 0,
			b.MQTT().WrapSubscribeDeviceIsOnline(func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.callbackMQTTTemperature(message, offset)
			})))
	}

	return b.MQTT().Subscribe(subscribers...)
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	cfg := b.config()

	stateParams := printer.NewGetPrinterStateParamsWithContext(ctx).
		WithHistory(&[]bool{false}[0]).
		WithExclude([]string{"sd"})

	if b.TemperatureFromMQTT() {
		stateParams = stateParams.WithExclude(append(stateParams.Exclude, "temperature"))
	}

	state, err := b.provider.Printer.GetPrinterState(stateParams, nil)
	id := b.Meta().ID()

	if err != nil {
		if _, ok := err.(*printer.GetPrinterStateConflict); ok {
			err = b.MQTT().PublishAsync(ctx, cfg.TopicState.Format(id), "Not operational")
		}

		return err
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicState.Format(id), state.Payload.State.Text); e != nil {
		err = multierr.Append(err, e)
	}

	if !b.TemperatureFromMQTT() {
		for device, temperature := range state.Payload.Temperature {
			metricDeviceTemperatureActual.With("id", id).With("device", device).Set(temperature.Actual)
			metricDeviceTemperatureTarget.With("id", id).With("device", device).Set(temperature.Target)

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicTemperatureActual.Format(id, device), temperature.Actual); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicTemperatureOffset.Format(id, device), temperature.Offset); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicTemperatureTarget.Format(id, device), temperature.Target); e != nil {
				err = multierr.Append(err, e)
			}
		}

		b.devicesMutex.RLock()
		count := len(b.devices)
		b.devicesMutex.RUnlock()

		if count == 0 {
			b.devicesMutex.Lock()
			for device := range state.Payload.Temperature {
				b.devices[device] = true
			}
			b.devicesMutex.Unlock()
		}
	}

	// Job
	if !b.JobFromMQTT() {
		if j, e := b.provider.Job.GetJob(job.NewGetJobParamsWithContext(ctx), nil); e == nil {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicJobFileName.Format(id), j.Payload.Job.File.Name); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicJobFileSize.Format(id), j.Payload.Job.File.Size); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicJobProgress.Format(id), j.Payload.Progress.Completion); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicJobTime.Format(id), j.Payload.Progress.PrintTime); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicJobTimeLeft.Format(id), j.Payload.Progress.PrintTimeLeft); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	// Layer & Height
	var (
		layerTotal, layerCurrent   uint64
		heightTotal, heightCurrent float64
	)

	if b.DisplayLayerProgressEnabled() && state.Payload.State.Flags.Printing {
		if progress, e := b.provider.Plugin.DisplayLayerProgress(plugin.NewDisplayLayerProgressParamsWithContext(ctx), nil); e == nil {
			if value, e := strconv.ParseUint(progress.Payload.Layer.Total, 10, 64); e == nil {
				layerTotal = value
			}

			if progress.Payload.Layer.Current != "-" {
				if value, e := strconv.ParseUint(progress.Payload.Layer.Current, 10, 64); e == nil {
					layerCurrent = value
				}
			}

			if value, e := strconv.ParseFloat(progress.Payload.Height.Total, 64); e == nil {
				heightTotal = value
			}

			if value, e := strconv.ParseFloat(progress.Payload.Height.Current, 64); e == nil {
				heightCurrent = value
			}
		}
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicLayerTotal.Format(id), layerTotal); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicLayerCurrent.Format(id), layerCurrent); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicHeightTotal.Format(id), heightTotal); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicHeightCurrent.Format(id), heightCurrent); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
