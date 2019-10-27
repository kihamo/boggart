package octoprint

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/octoprint/client/job"
	"github.com/kihamo/boggart/providers/octoprint/client/printer"
	"github.com/kihamo/boggart/providers/octoprint/models"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	/*
		loginParams := authorization.NewLoginParamsWithContext(ctx)
		loginParams.Body.Passive = true

		_, err := b.provider.Authorization.Login(loginParams, nil)
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
			return nil, err
		}
	*/

	stateParams := printer.NewGetPrinterStateParamsWithContext(ctx).
		WithHistory(&[]bool{false}[0]).
		WithExclude([]string{"sd"})

	state, err := b.provider.Printer.GetPrinterState(stateParams, nil)
	address := b.config.Address.Host

	if err != nil {
		if _, ok := err.(*printer.GetPrinterStateConflict); ok {
			err = b.MQTTPublishAsync(ctx, b.config.TopicState.Format(address), "Not operational")
			b.UpdateStatus(boggart.BindStatusOnline)
		} else {
			b.UpdateStatus(boggart.BindStatusOffline)
		}

		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	if e := b.MQTTPublishAsync(ctx, b.config.TopicState.Format(address), state.Payload.State.Text); e != nil {
		err = multierr.Append(err, e)
	}

	var temperature *models.TemperatureData

	// Bed
	temperature = state.Payload.Temperature.Bed

	metricDeviceTemperatureActual.With("address", address).With("device", "bed").Set(temperature.Actual)
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateBedTemperatureActual.Format(address), temperature.Actual); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateBedTemperatureOffset.Format(address), temperature.Offset); e != nil {
		err = multierr.Append(err, e)
	}

	metricDeviceTemperatureTarget.With("address", address).With("device", "bed").Set(temperature.Target)
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateBedTemperatureTarget.Format(address), temperature.Target); e != nil {
		err = multierr.Append(err, e)
	}

	// Hotend (tool 0)
	temperature = state.Payload.Temperature.Tool0

	metricDeviceTemperatureActual.With("address", address).With("device", "tool0").Set(temperature.Actual)
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateTool0TemperatureActual.Format(address), temperature.Actual); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateTool0TemperatureOffset.Format(address), temperature.Offset); e != nil {
		err = multierr.Append(err, e)
	}

	metricDeviceTemperatureTarget.With("address", address).With("device", "tool0").Set(temperature.Target)
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateTool0TemperatureTarget.Format(address), temperature.Target); e != nil {
		err = multierr.Append(err, e)
	}

	// Job
	if j, e := b.provider.Job.GetJob(job.NewGetJobParamsWithContext(ctx), nil); e == nil {
		if e := b.MQTTPublishAsync(ctx, b.config.TopicStateJobFileName.Format(address), j.Payload.Job.File.Name); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicStateJobFileSize.Format(address), j.Payload.Job.File.Size); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicStateJobProgress.Format(address), j.Payload.Progress.Completion); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicStateJobTime.Format(address), j.Payload.Progress.PrintTime); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicStateJobTimeLeft.Format(address), j.Payload.Progress.PrintTimeLeft); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
