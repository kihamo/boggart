package led_wifi

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
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
	var result error

	prevPower := atomic.LoadInt64(&b.statePower)
	if prevPower == 0 || (prevPower == 1) != state.Power {
		if state.Power {
			atomic.StoreInt64(&b.statePower, 1)
		} else {
			atomic.StoreInt64(&b.statePower, -1)
		}

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStatePower.Format(host), 0, true, state.Power); err != nil {
			result = multierr.Append(result, err)
		}
	}

	currentMode := uint64(state.Mode)
	prevMode := atomic.LoadUint64(&b.stateMode)
	if prevMode != currentMode {
		atomic.StoreUint64(&b.stateMode, currentMode)

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMode.Format(host), 0, true, currentMode); err != nil {
			result = multierr.Append(result, err)
		}
	}

	currentSpeed := uint64(state.Speed)
	prevSpeed := atomic.LoadUint64(&b.stateSpeed)
	if prevSpeed != currentSpeed {
		atomic.StoreUint64(&b.stateSpeed, currentSpeed)

		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateSpeed.Format(host), 0, true, currentSpeed); err != nil {
			result = multierr.Append(result, err)
		}
	}

	currentColor := state.Color.Uint64()
	prevColor := atomic.LoadUint64(&b.stateColor)
	if prevColor != currentColor {
		atomic.StoreUint64(&b.stateColor, currentColor)

		// in HEX
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColor.Format(host), 0, true, state.Color.String()); err != nil {
			result = multierr.Append(result, err)
		}

		// in HSV
		h, s, v := state.Color.HSV()
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColorHSV.Format(host), 0, true, fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); err != nil {
			result = multierr.Append(result, err)
		}
	}

	return nil, result
}
