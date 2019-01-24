package led_wifi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

const (
	TaskUpdaterInterval = time.Second * 3
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(TaskUpdaterInterval)
	taskUpdater.SetName("bind-led-wifi-updater-" + b.bulb.Host())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	state, err := b.bulb.State(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.SetSerialNumber(strconv.FormatUint(uint64(state.DeviceName), 10))

	b.UpdateStatus(boggart.BindStatusOnline)
	host := mqtt.NameReplace(b.bulb.Host())

	if ok := b.power.Set(state.Power); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStatePower.Format(host), 0, true, state.Power); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.mode.Set(uint32(state.Mode)); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMode.Format(host), 0, true, state.Mode); e != nil {
			err = multierr.Append(err, err)
		}
	}

	if ok := b.speed.Set(uint32(state.Speed)); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateSpeed.Format(host), 0, true, state.Speed); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.color.Set(uint32(state.Color.Uint64())); ok {
		// in HEX
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColor.Format(host), 0, true, state.Color.String()); e != nil {
			err = multierr.Append(err, e)
		}

		// in HSV
		h, s, v := state.Color.HSV()
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColorHSV.Format(host), 0, true, fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return nil, err
}
