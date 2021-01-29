package xmeye

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/xmeye"
	"go.uber.org/multierr"
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
	client, err := b.client(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	info, err := client.SystemInfo(ctx)
	if err != nil {
		return err
	}

	if info.SerialNo == "" {
		return errors.New("device returns empty serial number")
	}

	b.Meta().SetSerialNumber(info.SerialNo)

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
						b.config.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	)
	if err != nil {
		return err
	}

	if b.Meta().MAC() == nil {
		response, err := client.ConfigGet(ctx, xmeye.ConfigNameNetworkNetCommon, false)
		if err != nil {
			return err
		}

		if cfg, ok := response.(map[string]interface{}); ok {
			if mac, ok := cfg["MAC"]; ok {
				if err := b.Meta().SetMACAsString(mac.(string)); err != nil {
					return err
				}
			}
		}

		if b.Meta().MAC() == nil {
			return errors.New("device returns empty MAC address")
		}
	}

	if b.config.AlarmStreamingEnabled {
		if err = b.startAlarmStreaming(); err != nil {
			return err
		}
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(info.SerialNo), info.HardWare); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(info.SerialNo), info.SoftWareVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareReleasedDate.Format(info.SerialNo), info.BuildTime); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateUpTime.Format(info.SerialNo), time.Duration(info.DeviceRunTime)*time.Minute); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	client, err := b.client(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	sn := b.Meta().SerialNumber()

	storage, _ := client.StorageInfo(ctx)
	for _, s := range storage {
		for _, p := range s.Partition {
			if p.IsCurrent {
				name := strconv.FormatUint(p.LogicSerialNo, 10)

				metricStorageUsage.With("serial_number", sn).With("name", name).Set(float64(uint64(p.TotalSpace-p.RemainSpace) * MB))
				metricStorageAvailable.With("serial_number", sn).With("name", name).Set(float64(uint64(p.RemainSpace) * MB))

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDCapacity.Format(sn, p.LogicSerialNo), uint64(p.TotalSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDUsage.Format(sn, p.LogicSerialNo), uint64(p.TotalSpace-p.RemainSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDFree.Format(sn, p.LogicSerialNo), uint64(p.RemainSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}
			}
		}
	}

	return err
}
