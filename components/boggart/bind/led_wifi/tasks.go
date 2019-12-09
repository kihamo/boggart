package led_wifi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater-" + b.config.Address)

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, nil
	}

	state, err := b.bulb.State(ctx)
	if err != nil {
		return nil, err
	}

	b.SetSerialNumber(strconv.FormatUint(uint64(state.DeviceName), 10))

	if e := b.MQTTPublishAsync(ctx, b.config.TopicStatePower, state.Power); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateMode, state.Mode); e != nil {
		err = multierr.Append(err, err)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateSpeed, state.Speed); e != nil {
		err = multierr.Append(err, e)
	}

	// in HEX
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateColor, state.Color.String()); e != nil {
		err = multierr.Append(err, e)
	}

	// in HSV
	h, s, v := state.Color.HSV()
	if e := b.MQTTPublishAsync(ctx, b.config.TopicStateColorHSV, fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
		err = multierr.Append(err, e)
	}

	return nil, err
}
