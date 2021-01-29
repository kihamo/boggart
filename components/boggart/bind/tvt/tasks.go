package tvt

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/tvt/client/information"
	"github.com/kihamo/boggart/providers/tvt/client/net"
	"github.com/kihamo/boggart/providers/tvt/client/storage"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
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
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		baseCfg, e := b.client.Information.GetBasicConfig(information.NewGetBasicConfigParamsWithContext(ctx), nil)
		if e != nil {
			return e
		}

		if baseCfg.Payload.Content.Sn == "" {
			return errors.New("device returns empty serial number")
		}

		b.Meta().SetSerialNumber(baseCfg.Payload.Content.Sn)
		sn = baseCfg.Payload.Content.Sn

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(sn), baseCfg.Payload.Content.ProductModel); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(sn), baseCfg.Payload.Content.SoftwareVersion); e != nil {
			err = multierr.Append(err, e)
		}
	}

	mac := b.Meta().MAC()
	if mac == nil {
		netStatus, e := b.client.Net.GetNetStatus(net.NewGetNetStatusParamsWithContext(ctx), nil)
		if e != nil {
			return e
		}

		primaryNic := netStatus.Payload.Content.IPGroup.PrimaryNIC
		if primaryNic == "" {
			return errors.New("device returns empty primary nic")
		}

		for _, nic := range netStatus.Payload.Content.Nic {
			if nic.ID == primaryNic {
				if e := b.Meta().SetMACAsString(nic.Mac); e != nil {
					return e
				}

				break
			}
		}
	}

	disks, e := b.client.Storage.GetStorageInfo(storage.NewGetStorageInfoParamsWithContext(ctx), nil)
	if e != nil {
		return e
	}

	var unitFactor int64 = 1

	if disks.Payload.Content.DisksSize.Unit == "MB" {
		unitFactor = 1024
	}

	for _, disk := range disks.Payload.Content.Disks {
		capacity := float64(disk.Size * unitFactor)
		free := float64(disk.FreeSpace * unitFactor)
		usage := capacity - free

		metricStorageUsage.With("serial_number", sn).With("name", disk.SerialNum).Set(usage)
		metricStorageAvailable.With("serial_number", sn).With("name", disk.SerialNum).Set(free)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDCapacity.Format(sn, disk.SerialNum), capacity); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDUsage.Format(sn, disk.SerialNum), usage); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDFree.Format(sn, disk.SerialNum), free); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
