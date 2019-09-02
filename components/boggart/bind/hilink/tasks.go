package hilink

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/monitoring"
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

	taskBalanceUpdater := task.NewFunctionTask(b.taskBalanceUpdater)
	taskBalanceUpdater.SetTimeout(b.config.BalanceUpdaterTimeout)
	taskBalanceUpdater.SetRepeats(-1)
	taskBalanceUpdater.SetRepeatInterval(b.config.BalanceUpdaterInterval)
	taskBalanceUpdater.SetName("bind-hilink-balance-updater-" + b.config.Address.Host)

	taskSMSChecker := task.NewFunctionTask(b.taskSMSChecker)
	taskSMSChecker.SetTimeout(b.config.SMSCheckerTimeout)
	taskSMSChecker.SetRepeats(-1)
	taskSMSChecker.SetRepeatInterval(b.config.SMSCheckerInterval)
	taskSMSChecker.SetName("bind-hilink-sms-checker-" + b.config.Address.Host)

	taskSignalUpdater := task.NewFunctionTask(b.taskSignalUpdater)
	taskSignalUpdater.SetTimeout(b.config.SignalUpdaterTimeout)
	taskSignalUpdater.SetRepeats(-1)
	taskSignalUpdater.SetRepeatInterval(b.config.SignalUpdaterInterval)
	taskSignalUpdater.SetName("bind-hilink-signal-updater-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskBalanceUpdater,
		taskSMSChecker,
		taskSignalUpdater,
	}

	return tasks
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	deviceInfo, err := b.client.Device.GetDeviceInformation(device.NewGetDeviceInformationParamsWithContext(ctx))
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	if deviceInfo.Payload.SerialNumber == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("device returns empty serial number")
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(deviceInfo.Payload.SerialNumber)
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	if b.operator.IsEmpty() {
		plmn, err := b.client.Net.GetCurrentPLMN(net.NewGetCurrentPLMNParamsWithContext(ctx))
		if err == nil {
			b.operator.Set(plmn.Payload.FullName)
			return nil, b.MQTTPublishAsync(
				ctx,
				MQTTPublishTopicOperator.Format(mqtt.NameReplace(deviceInfo.Payload.SerialNumber)),
				plmn.Payload.FullName)
		} else {
			return nil, err
		}
	}

	return nil, err
}

func (b *Bind) taskBalanceUpdater(ctx context.Context) (interface{}, error) {
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
				if e = b.MQTTPublishAsync(ctx, MQTTPublishTopicSMS.Format(sn), payload); e != nil {
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

func (b *Bind) taskSignalUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	paramsSignal := device.NewGetDeviceSignalParamsWithContext(ctx)
	responseSignal, err := b.client.Device.GetDeviceSignal(paramsSignal)
	if err != nil {
		return nil, err
	}

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)

	if val, e := strconv.ParseInt(strings.TrimRight(responseSignal.Payload.RSSI, "dBm"), 10, 64); e == nil {
		metricSignalRSSI.With("serial_number", sn).Set(float64(val))

		if e = b.MQTTPublishAsync(ctx, MQTTPublishSignalRSSI.Format(snMQTT), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if val, e := strconv.ParseInt(strings.TrimRight(responseSignal.Payload.RSRP, "dBm"), 10, 64); e == nil {
		metricSignalRSRP.With("serial_number", sn).Set(float64(val))

		if e = b.MQTTPublishAsync(ctx, MQTTPublishSignalRSRP.Format(snMQTT), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if val, e := strconv.ParseInt(strings.TrimRight(responseSignal.Payload.RSRQ, "dBm"), 10, 64); e == nil {
		metricSignalRSRQ.With("serial_number", sn).Set(float64(val))

		if e = b.MQTTPublishAsync(ctx, MQTTPublishSignalRSRQ.Format(snMQTT), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	if val, e := strconv.ParseInt(strings.TrimRight(responseSignal.Payload.SINR, "dBm"), 10, 64); e == nil {
		metricSignalSINR.With("serial_number", sn).Set(float64(val))

		if e = b.MQTTPublishAsync(ctx, MQTTPublishSignalSINR.Format(snMQTT), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	paramsStatus := monitoring.NewGetMonitoringStatusParamsWithContext(ctx)
	responseStatus, err := b.client.Monitoring.GetMonitoringStatus(paramsStatus)
	if err != nil {
		return nil, err
	}

	val := responseStatus.Payload.SignalIcon
	metricSignalLevel.With("serial_number", sn).Set(float64(val))

	if e := b.MQTTPublishAsync(ctx, MQTTPublishSignalLevel.Format(snMQTT), val); e != nil {
		err = multierr.Append(err, e)
	}

	return nil, err
}
