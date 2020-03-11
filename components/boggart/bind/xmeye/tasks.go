package xmeye

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/xmeye"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := b.Workers().WrapTaskIsOnlineOnceSuccess(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) error {
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

	taskState := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("updater")
	b.Workers().RegisterTask(taskState)

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

	return err
}

func (b *Bind) taskUpdater(ctx context.Context) (err error) {
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

				// TODO:
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDCapacity.Format(sn, p.LogicSerialNo), uint64(p.TotalSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}

				// TODO:
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDUsage.Format(sn, p.LogicSerialNo), uint64(p.TotalSpace-p.RemainSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}
				metricStorageUsage.With("serial_number", sn).With("name", name).Set(float64(uint64(p.TotalSpace-p.RemainSpace) * MB))

				// TODO:
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDFree.Format(sn, p.LogicSerialNo), uint64(p.RemainSpace)*MB); e != nil {
					err = multierr.Append(err, e)
				}
				metricStorageAvailable.With("serial_number", sn).With("name", name).Set(float64(uint64(p.RemainSpace) * MB))
			}
		}
	}

	return err
}
