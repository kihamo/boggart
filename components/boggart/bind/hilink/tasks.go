package hilink

import (
	"context"
	"errors"
	"github.com/kihamo/boggart/components/mqtt"
	"go.uber.org/multierr"

	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/sms"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/static/models"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/device"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-hilink-liveness-" + b.config.Address.Host)

	taskSMSChecker := task.NewFunctionTask(b.taskSMSChecker)
	taskSMSChecker.SetTimeout(b.config.LivenessTimeout)
	taskSMSChecker.SetRepeats(-1)
	taskSMSChecker.SetRepeatInterval(b.config.LivenessInterval)
	taskSMSChecker.SetName("bind-hilink-sms-checker-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskSMSChecker,
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := b.client.Device.GetDeviceInformation(device.NewGetDeviceInformationParamsWithContext(ctx))
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if deviceInfo.Payload.SerialNumber == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(deviceInfo.Payload.SerialNumber)
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, err
}

func (b *Bind) taskSMSChecker(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	// ----- read new sms -----
	params := sms.NewGetSMSListParamsWithContext(ctx).
		WithRequest(&models.SMSListRequest{
			PageIndex: 1,
			ReadCount: 20,
			BoxType:   1,
		})

	response, err := b.client.Sms.GetSMSList(params)
	if err != nil {
		return nil, err
	}

	sn := mqtt.NameReplace(b.SerialNumber())

	for _, s := range response.Payload.Messages {
		if s.Status != 1 {
			e := b.MQTTPublishAsync(ctx, MQTTPublishTopicSMS.Format(sn), s.Content)
			if e != nil {
				err = multierr.Append(err, e)
				continue
			}

			params := sms.NewReadSMSParamsWithContext(ctx)
			params.Request.Index = s.Index

			_, e = b.client.Sms.ReadSMS(params)
			if e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	// ----- delete old sms -----

	return nil, err
}
