package nut

import (
	"context"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithDurationFunc(
					tasks.ScheduleNow(),
					func(meta tasks.Meta) time.Duration {
						return b.updaterInterval.Load()
					},
				),
			),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	variables, err := b.Variables()
	if err != nil {
		return err
	}

	sn := b.Meta().SerialNumber()

	if sn == "" {
		for _, v := range variables {
			if v.Name == "device.serial" {
				sn = strings.TrimSpace(v.Value.(string))
				b.Meta().SetSerialNumber(sn)

				break
			}
		}
	}

	if sn == "" {
		return nil
	}

	for _, v := range variables {
		switch v.Name {
		case "ups.load":
			metricLoad.With("serial_number", sn).Set(float64(v.Value.(int)))
		case "input.voltage":
			metricInputVoltage.With("serial_number", sn).Set(v.Value.(float64))
		case "battery.charge":
			metricBatteryCharge.With("serial_number", sn).Set(float64(v.Value.(int)))
		case "battery.runtime":
			metricBatteryRuntime.With("serial_number", sn).Set(float64(v.Value.(int)))
		case "battery.voltage":
			metricBatteryVoltage.With("serial_number", sn).Set(v.Value.(float64))
		case "driver.parameter.pollfreq":
			b.updaterInterval.Set(time.Duration(v.Value.(int)) * time.Second)
		}

		// TODO:
		_ = b.MQTT().PublishAsync(ctx, b.config.TopicVariable.Format(sn, v.Name), v.Value)
	}

	return nil
}
