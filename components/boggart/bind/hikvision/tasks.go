package hikvision

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/hikvision/client/content_manager"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
	"github.com/kihamo/boggart/providers/hikvision/client/system"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
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
	deviceInfo, err := b.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)
	if err != nil {
		return err
	}

	if deviceInfo.Payload.SerialNumber == "" {
		return errors.New("device returns empty serial number")
	}

	if deviceInfo.Payload.MacAddress == "" {
		return errors.New("device returns empty MAC address")
	}

	b.Meta().SetSerialNumber(deviceInfo.Payload.SerialNumber)

	err = b.Meta().SetMACAsString(deviceInfo.Payload.MacAddress)
	if err != nil {
		return err
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.Model); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareReleasedDate.Format(deviceInfo.Payload.SerialNumber), deviceInfo.Payload.FirmwareReleasedDate); e != nil {
		err = multierr.Append(err, e)
	}

	_, e := b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler), b.config.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	if b.config.EventsEnabled && b.config.EventsStreamingEnabled {
		b.startAlertStreaming()
	}

	if channels, e := b.client.Ptz.GetPtzChannels(ptz.NewGetPtzChannelsParamsWithContext(ctx), nil); e == nil && len(channels.Payload) > 0 {
		// FIXME: почему-то при ответе
		// <?xml version="1.0" encoding="UTF-8"?>
		// <PTZChannelList  version="2.0" xmlns="http://www.std-cgi.com/ver20/XMLSchema">
		// </PTZChannelList >
		// создается 1 пустой элемент, поэтому быстрохак
		if !(len(channels.Payload) == 1 && channels.Payload[0].ID == 0) {
			e = b.MQTT().Subscribe(
				mqtt.NewSubscriber(b.config.TopicPTZAbsolute, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTAbsolute)),
				mqtt.NewSubscriber(b.config.TopicPTZContinuous, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTContinuous)),
				mqtt.NewSubscriber(b.config.TopicPTZRelative, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTRelative)),
				mqtt.NewSubscriber(b.config.TopicPTZPreset, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTPreset)),
				mqtt.NewSubscriber(b.config.TopicPTZMomentary, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTMomentary)),
			)
			if e == nil {
				err = multierr.Append(err, e)
			}

			_, e = b.Workers().RegisterTask(
				tasks.NewTask().
					WithName("ptz").
					WithHandler(
						b.Workers().WrapTaskIsOnline(
							tasks.HandlerWithTimeout(
								tasks.HandlerFuncFromShortToLong(b.taskPTZHandler), b.config.PTZTimeout,
							),
						),
					).
					WithSchedule(
						tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.PTZInterval),
					),
			)
			if e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

func (b *Bind) taskPTZHandler(ctx context.Context) error {
	channels, err := b.client.Ptz.GetPtzChannels(ptz.NewGetPtzChannelsParamsWithContext(ctx), nil)
	if err != nil {
		return err
	}

	for _, ch := range channels.Payload {
		if !ch.Enabled {
			continue
		}

		if err := b.updateStatusByChannelID(ctx, ch.ID); err != nil {
			b.Logger().Error("Failed updated status",
				"serial_number", b.Meta().SerialNumber(),
				"channel", ch.ID,
				"error", err.Error(),
			)
		}
	}

	return nil
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	if status, e := b.client.System.GetStatus(system.NewGetStatusParamsWithContext(ctx), nil); e == nil {
		memoryUsage := int64(status.Payload.MemoryList[0].MemoryUsage) * MB
		memoryAvailable := int64(status.Payload.MemoryList[0].MemoryAvailable) * MB

		metricUpTime.With("serial_number", sn).Set(float64(status.Payload.DeviceUpTime))
		metricMemoryUsage.With("serial_number", sn).Set(float64(memoryUsage))
		metricMemoryAvailable.With("serial_number", sn).Set(float64(memoryAvailable))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateUpTime.Format(sn), status.Payload.DeviceUpTime); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateMemoryUsage.Format(sn), memoryUsage); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateMemoryAvailable.Format(sn), memoryAvailable); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if storage, e := b.client.ContentManager.GetStorage(content_manager.NewGetStorageParamsWithContext(ctx), nil); e == nil {
		for _, hdd := range storage.Payload.HddList {
			if hdd.Name == "" {
				continue
			}

			metricStorageUsage.With("serial_number", sn).With("name", hdd.Name).Set(float64((hdd.Capacity - hdd.FreeSpace) * MB))
			metricStorageAvailable.With("serial_number", sn).With("name", hdd.Name).Set(float64(hdd.FreeSpace * MB))

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDCapacity.Format(sn, hdd.ID), hdd.Capacity*MB); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDUsage.Format(sn, hdd.ID), (hdd.Capacity-hdd.FreeSpace)*MB); e != nil {
				err = multierr.Append(err, e)
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateHDDFree.Format(sn, hdd.ID), hdd.FreeSpace*MB); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}
