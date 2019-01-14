package led_wifi

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/mqtt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
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

	prevPower := atomic.LoadInt64(&b.statePower)
	if prevPower == 0 || (prevPower == 1) != state.Power {
		if state.Power {
			atomic.StoreInt64(&b.statePower, 1)
		} else {
			atomic.StoreInt64(&b.statePower, -1)
		}

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStatePower.Format(host), 0, true, state.Power)
	}

	currentMode := uint64(state.Mode)
	prevMode := atomic.LoadUint64(&b.stateMode)
	if prevMode != currentMode {
		atomic.StoreUint64(&b.stateMode, currentMode)

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMode.Format(host), 0, true, currentMode)
	}

	currentSpeed := uint64(state.Speed)
	prevSpeed := atomic.LoadUint64(&b.stateSpeed)
	if prevSpeed != currentSpeed {
		atomic.StoreUint64(&b.stateSpeed, currentSpeed)

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateSpeed.Format(host), 0, true, currentSpeed)
	}

	currentColor := state.Color.Uint64()
	prevColor := atomic.LoadUint64(&b.stateColor)
	if prevColor != currentColor {
		atomic.StoreUint64(&b.stateColor, currentColor)

		// in HEX
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColor.Format(host), 0, true, state.Color.String())

		// in HSV
		h, s, v := state.Color.HSV()
		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateColorHSV.Format(host), 0, true, fmt.Sprintf("%d,%.2f,%.2f", h, s, v))
	}

	return nil, nil
}
