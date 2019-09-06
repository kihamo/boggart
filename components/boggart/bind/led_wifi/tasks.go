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
	taskUpdater.SetName("updater-" + b.bulb.Host())

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

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStatePower.Format(host), state.Power); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMode.Format(host), state.Mode); e != nil {
		err = multierr.Append(err, err)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateSpeed.Format(host), state.Speed); e != nil {
		err = multierr.Append(err, e)
	}

	// in HEX
	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColor.Format(host), state.Color.String()); e != nil {
		err = multierr.Append(err, e)
	}

	// in HSV
	h, s, v := state.Color.HSV()
	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColorHSV.Format(host), fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
		err = multierr.Append(err, e)
	}

	return nil, err
}
