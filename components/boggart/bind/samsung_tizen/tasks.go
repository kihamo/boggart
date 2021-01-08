package tizen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumber),
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

func (b *Bind) taskSerialNumber(ctx context.Context) error {
	info, err := b.client.Device(ctx)
	if err != nil {
		return err
	}

	parts := strings.Split(info.ID, ":")
	if len(parts) > 1 {
		b.Meta().SetSerialNumber(parts[1])
	} else {
		b.Meta().SetSerialNumber(info.ID)
	}

	if e := b.Meta().SetMACAsString(info.Device.WifiMac); e != nil {
		err = fmt.Errorf("set mac address failed: %w", err)
	}

	mac := b.Meta().MACAsString()

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicDeviceID.Format(mac), info.Device.ID); e != nil {
		err = fmt.Errorf("send mqtt message about device id failed: %w", err)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicDeviceModelName.Format(mac), info.Device.Name); e != nil {
		err = fmt.Errorf("send mqtt message about device model name failed: %w", err)
	}

	return err
}
