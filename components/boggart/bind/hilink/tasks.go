package hilink

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/providers/hilink/client/monitoring"
	"github.com/kihamo/boggart/providers/hilink/client/net"
	"github.com/kihamo/boggart/providers/hilink/client/sms"
	"github.com/kihamo/boggart/providers/hilink/static/models"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.config.Address.Host)

	taskBalanceUpdater := task.NewFunctionTask(b.taskBalanceUpdater)
	taskBalanceUpdater.SetTimeout(b.config.BalanceUpdaterTimeout)
	taskBalanceUpdater.SetRepeats(-1)
	taskBalanceUpdater.SetRepeatInterval(b.config.BalanceUpdaterInterval)
	taskBalanceUpdater.SetName("balance-updater-" + b.config.Address.Host)

	taskSMSChecker := task.NewFunctionTask(b.taskSMSChecker)
	taskSMSChecker.SetTimeout(b.config.SMSCheckerTimeout)
	taskSMSChecker.SetRepeats(-1)
	taskSMSChecker.SetRepeatInterval(b.config.SMSCheckerInterval)
	taskSMSChecker.SetName("sms-checker-" + b.config.Address.Host)

	taskSystemUpdater := task.NewFunctionTask(b.taskSystemUpdater)
	taskSystemUpdater.SetTimeout(b.config.SystemUpdaterTimeout)
	taskSystemUpdater.SetRepeats(-1)
	taskSystemUpdater.SetRepeatInterval(b.config.SystemUpdaterInterval)
	taskSystemUpdater.SetName("system-updater-" + b.config.Address.Host)

	tasks := []workers.Task{
		taskLiveness,
		taskBalanceUpdater,
		taskSMSChecker,
		taskSystemUpdater,
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
				b.config.TopicOperator.Format(deviceInfo.Payload.SerialNumber),
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

	balance, err := b.Balance(ctx)
	if err == nil {
		metricBalance.With("serial_number", sn).Set(balance)

		if e := b.MQTTPublishAsync(ctx, b.config.TopicBalance.Format(sn), balance); e != nil {
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

	sn := b.SerialNumber()

	for _, s := range response.Payload.Messages {
		if s.Status != 1 {
			payload, e := json.Marshal(s)
			if err != nil {
				err = multierr.Append(err, e)
				continue
			}

			isSpecial := b.checkSpecialSMS(ctx, s)
			if !isSpecial {
				if e = b.MQTTPublishAsync(ctx, b.config.TopicSMS.Format(sn), payload); e != nil {
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

func (b *Bind) taskSystemUpdater(ctx context.Context) (_ interface{}, err error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	sn := b.SerialNumber()

	// signal
	if response, e := b.client.Device.GetDeviceSignal(device.NewGetDeviceSignalParamsWithContext(ctx)); e == nil {
		if val, e := strconv.ParseInt(strings.TrimRight(response.Payload.RSSI, "dBm"), 10, 64); e == nil {
			metricSignalRSSI.With("serial_number", sn).Set(float64(val))

			if e = b.MQTTPublishAsync(ctx, b.config.TopicSignalRSSI.Format(sn), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}

		if val, e := strconv.ParseInt(strings.TrimRight(response.Payload.RSRP, "dBm"), 10, 64); e == nil {
			metricSignalRSRP.With("serial_number", sn).Set(float64(val))

			if e = b.MQTTPublishAsync(ctx, b.config.TopicSignalRSRP.Format(sn), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}

		if val, e := strconv.ParseInt(strings.TrimRight(response.Payload.RSRQ, "dBm"), 10, 64); e == nil {
			metricSignalRSRQ.With("serial_number", sn).Set(float64(val))

			if e = b.MQTTPublishAsync(ctx, b.config.TopicSignalRSRQ.Format(sn), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}

		if val, e := strconv.ParseInt(strings.TrimRight(response.Payload.SINR, "dBm"), 10, 64); e == nil {
			metricSignalSINR.With("serial_number", sn).Set(float64(val))

			if e = b.MQTTPublishAsync(ctx, b.config.TopicSignalSINR.Format(sn), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	// status
	if response, e := b.client.Monitoring.GetMonitoringStatus(monitoring.NewGetMonitoringStatusParamsWithContext(ctx)); e == nil {
		metricSignalLevel.With("serial_number", sn).Set(float64(response.Payload.SignalIcon))
		if e := b.MQTTPublishAsync(ctx, b.config.TopicSignalLevel.Format(sn), response.Payload.SignalIcon); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	// traffic
	if response, e := b.client.Monitoring.GetMonitoringTrafficStatistics(monitoring.NewGetMonitoringTrafficStatisticsParamsWithContext(ctx)); e == nil {
		metricTotalConnectTime.With("serial_number", sn).Set(float64(response.Payload.TotalConnectTime))
		if e := b.MQTTPublishAsync(ctx, b.config.TopicConnectionTime.Format(sn), response.Payload.TotalConnectTime); e != nil {
			err = multierr.Append(err, e)
		}

		metricTotalDownload.With("serial_number", sn).Set(float64(response.Payload.TotalDownload))
		if e := b.MQTTPublishAsync(ctx, b.config.TopicConnectionDownload.Format(sn), response.Payload.TotalDownload); e != nil {
			err = multierr.Append(err, e)
		}

		metricTotalUpload.With("serial_number", sn).Set(float64(response.Payload.TotalUpload))
		if e := b.MQTTPublishAsync(ctx, b.config.TopicConnectionUpload.Format(sn), response.Payload.TotalUpload); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}
