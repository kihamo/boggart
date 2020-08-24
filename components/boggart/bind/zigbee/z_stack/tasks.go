package zstack

import (
	"context"
	"encoding/hex"
	"time"

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

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionTransportRevision.Format(sn), version.TransportRevision); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionProduct.Format(sn), version.Product); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionMajorRelease.Format(sn), version.MajorRelease); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionMinorRelease.Format(sn), version.MinorRelease); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionMainTrel.Format(sn), version.MainTrel); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionHardwareRevision.Format(sn), version.HardwareRevision); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersionType.Format(sn), version.Type()); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
