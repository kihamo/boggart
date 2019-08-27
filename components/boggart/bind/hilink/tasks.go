package hilink

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/net"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/sms"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/static/models"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-hilink-liveness-" + b.config.Address.Host)

	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetTimeout(b.config.UpdaterTimeout)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("bind-hilink-updater-" + b.config.Address.Host)

	taskSMSChecker := task.NewFunctionTask(b.taskSMSChecker)
	taskSMSChecker.SetTimeout(b.config.SMSCheckerTimeout)
	taskSMSChecker.SetRepeats(-1)
	taskSMSChecker.SetRepeatInterval(b.config.SMSCheckerInterval)
	taskSMSChecker.SetName("bind-hilink-sms-checker-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskUpdater,
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

	if b.operator.IsEmpty() {
		plmn, err := b.client.Net.GetCurrentPLMN(net.NewGetCurrentPLMNParamsWithContext(ctx))
		if err == nil {
			b.operator.Set(plmn.Payload.FullName)
			b.MQTTPublishAsync(
				ctx,
				MQTTPublishTopicOperator.Format(mqtt.NameReplace(deviceInfo.Payload.SerialNumber)),
				plmn.Payload.FullName)
		}
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, err
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)

	balance, err := b.Balance(ctx)
	if err == nil {
		metricBalance.With("serial_number", sn).Set(balance)

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBalance.Format(snMQTT), balance); e != nil {
			err = multierr.Append(err, e)
		}
	}

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
			payload, e := json.Marshal(s)
			if err != nil {
				err = multierr.Append(err, e)
				continue
			}

			isSpecial := b.checkSpecialSMS(ctx, s)
			if !isSpecial {
				e = b.MQTTPublishAsync(ctx, MQTTPublishTopicSMS.Format(sn), payload)
				if e != nil {
					err = multierr.Append(err, e)
					continue
				}
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
