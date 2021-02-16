package zstack

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
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
	client, err := b.getClient()
	if err != nil {
		return err
	}

	info, err := client.UtilGetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	sn := hex.EncodeToString(info.IEEEAddr)
	b.Meta().SetSerialNumber(sn)
	b.syncPermitJoin()

	version, err := client.Version(ctx)
	if err != nil {
		return err
	}

	cfg := b.config()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionTransportRevision.Format(sn), version.TransportRevision); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionProduct.Format(sn), version.Product); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionMajorRelease.Format(sn), version.MajorRelease); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionMinorRelease.Format(sn), version.MinorRelease); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionMainTrel.Format(sn), version.MainTrel); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionHardwareRevision.Format(sn), version.HardwareRevision); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersionType.Format(sn), version.Type()); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
