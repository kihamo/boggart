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
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumberHandler),
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

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
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
	cfg := b.config()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicDeviceID.Format(mac), info.Device.ID); e != nil {
		err = fmt.Errorf("send mqtt message about device id failed: %w", err)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicDeviceModelName.Format(mac), info.Device.Name); e != nil {
		err = fmt.Errorf("send mqtt message about device model name failed: %w", err)
	}

	return err
}
