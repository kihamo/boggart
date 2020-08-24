package hilink

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/providers/hilink/client/global"
	"github.com/kihamo/boggart/providers/hilink/client/monitoring"
	"github.com/kihamo/boggart/providers/hilink/client/net"
	"github.com/kihamo/boggart/providers/hilink/client/sms"
	"github.com/kihamo/boggart/providers/hilink/static/models"
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
	deviceInfo, err := b.client.Device.GetDeviceInformation(device.NewGetDeviceInformationParamsWithContext(ctx))
	if err != nil {
		return err
	}

	if deviceInfo.Payload.SerialNumber == "" {
		return errors.New("device returns empty serial number")
	}

	if deviceInfo.Payload.MacAddress1 == "" && deviceInfo.Payload.MacAddress2 == "" {
		return errors.New("device returns empty MAC address")
	}

	if deviceInfo.Payload.MacAddress1 != "" {
		err = b.Meta().SetMACAsString(deviceInfo.Payload.MacAddress1)
	} else if deviceInfo.Payload.MacAddress2 != "" {
		err = b.Meta().SetMACAsString(deviceInfo.Payload.MacAddress2)
	}

	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(deviceInfo.Payload.SerialNumber)

	// set settings
	settings, err := b.client.Global.GetGlobalModuleSwitch(global.NewGetGlobalModuleSwitchParamsWithContext(ctx))
	if err != nil {
		return err
	}

	ussdEnabled := settings.Payload.USSDEnabled > 0
	smsEnabled := settings.Payload.SMSEnabled > 0

	if ussdEnabled {
		taskBalanceUpdater := b.Workers().WrapTaskIsOnline(b.taskBalanceUpdater)
		taskBalanceUpdater.SetTimeout(b.config.BalanceUpdaterTimeout)
		taskBalanceUpdater.SetRepeats(-1)
		taskBalanceUpdater.SetRepeatInterval(b.config.BalanceUpdaterInterval)
		taskBalanceUpdater.SetName("balance-updater")
		b.Workers().RegisterTask(taskBalanceUpdater)

		taskLimitTrafficUpdater := b.Workers().WrapTaskIsOnline(b.taskLimitTrafficUpdater)
		taskLimitTrafficUpdater.SetTimeout(b.config.LimitTrafficUpdaterTimeout)
		taskLimitTrafficUpdater.SetRepeats(-1)
		taskLimitTrafficUpdater.SetRepeatInterval(b.config.LimitTrafficUpdaterInterval)
		taskLimitTrafficUpdater.SetName("limit-traffic-updater")
		b.Workers().RegisterTask(taskLimitTrafficUpdater)
	}

	if smsEnabled {
		taskSMSChecker := b.Workers().WrapTaskIsOnline(b.taskSMSChecker)
		taskSMSChecker.SetTimeout(b.config.SMSCheckerTimeout)
		taskSMSChecker.SetRepeats(-1)
		taskSMSChecker.SetRepeatInterval(b.config.SMSCheckerInterval)
		taskSMSChecker.SetName("sms-checker")
		b.Workers().RegisterTask(taskSMSChecker)

		taskCleaner := b.Workers().WrapTaskIsOnline(b.taskCleaner)
		taskCleaner.SetRepeats(-1)
		taskCleaner.SetRepeatInterval(b.config.CleanerInterval)
		taskCleaner.SetName("cleaner")
		b.Workers().RegisterTask(taskCleaner)
	}

	taskSystemUpdater := b.Workers().WrapTaskIsOnline(b.taskSystemUpdater)
	taskSystemUpdater.SetTimeout(b.config.SystemUpdaterTimeout)
	taskSystemUpdater.SetRepeats(-1)
	taskSystemUpdater.SetRepeatInterval(b.config.SystemUpdaterInterval)
	taskSystemUpdater.SetName("system-updater")
	b.Workers().RegisterTask(taskSystemUpdater)

	return nil
}

func (b *Bind) taskBalanceUpdater(ctx context.Context) error {
	if b.simStatus.Load() != 1 {
		return nil
	}

	balance, err := b.Balance(ctx)

	if err == nil {
		sn := b.Meta().SerialNumber()

		metricBalance.With("serial_number", sn).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBalance.Format(sn), balance); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

func (b *Bind) taskLimitTrafficUpdater(ctx context.Context) error {
	if b.simStatus.Load() != 1 || b.operator.IsEmpty() {
		return nil
	}

	op, err := b.operatorSettings()
	if err != nil {
		return err
	}

	if op.LimitTrafficUSSD == "" {
		return nil
	}

	_, err = b.USSD(ctx, op.LimitTrafficUSSD)

	return err
}

func (b *Bind) taskSMSChecker(ctx context.Context) error {
	if b.simStatus.Load() != 1 {
		return nil
	}

	// sms counters
	paramsCount := sms.NewGetSMSCountParamsWithContext(ctx)
	responseCount, err := b.client.Sms.GetSMSCount(paramsCount)

	if err != nil {
		return err
	}

	sn := b.Meta().SerialNumber()

	metricSMSUnread.With("serial_number", sn).Set(float64(responseCount.Payload.LocalUnread))
	metricSMSInbox.With("serial_number", sn).Set(float64(responseCount.Payload.LocalInbox))

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSMSUnread.Format(sn), responseCount.Payload.LocalUnread); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicSMSInbox.Format(sn), responseCount.Payload.LocalInbox); e != nil {
		err = multierr.Append(err, e)
	}

	// ----- read new sms -----
	paramsList := sms.NewGetSMSListParamsWithContext(ctx).
		WithRequest(&models.SMSListRequest{
			PageIndex: 1,
			ReadCount: 20,
			BoxType:   1,
		})

	responseList, e := b.client.Sms.GetSMSList(paramsList)
	if e != nil {
		return multierr.Append(err, e)
	}

	for _, s := range responseList.Payload.Messages {
		if s.Status != 1 {
			payload, e := json.Marshal(s)

			if err != nil {
				err = multierr.Append(err, e)
				continue
			}

			isSpecial := b.checkSpecialSMS(ctx, s)
			if !isSpecial {
				if e = b.MQTT().PublishAsync(ctx, b.config.TopicSMS.Format(sn), payload); e != nil {
					err = multierr.Append(err, e)
					continue
				}
			}

			if isSpecial && b.config.CleanerSpecial {
				params := sms.NewRemoveSMSParamsWithContext(ctx)
				params.Request.Index = s.Index

				_, e = b.client.Sms.RemoveSMS(params)
			} else {
				params := sms.NewReadSMSParamsWithContext(ctx)
				params.Request.Index = s.Index

				_, e = b.client.Sms.ReadSMS(params)
			}

			if e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	// ----- delete old sms -----

	return err
}

func (b *Bind) taskSystemUpdater(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	// status
	if response, e := b.client.Monitoring.GetMonitoringStatus(monitoring.NewGetMonitoringStatusParamsWithContext(ctx)); e == nil {
		b.simStatus.Set(uint32(response.Payload.SimStatus))

		// только с активной SIM картой
		if response.Payload.SimStatus == 1 {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicSignalLevel.Format(sn), response.Payload.SignalIcon); e != nil {
				err = multierr.Append(err, e)
			}

			metricSignalLevel.With("serial_number", sn).Set(float64(response.Payload.SignalIcon))

			if b.operator.IsEmpty() {
				plmn, e := b.client.Net.GetCurrentPLMN(net.NewGetCurrentPLMNParamsWithContext(ctx))
				if e == nil && plmn.Payload.FullName != "" {
					b.operator.Set(plmn.Payload.FullName)

					if e := b.MQTT().PublishAsync(ctx, b.config.TopicOperator.Format(sn), plmn.Payload.FullName); e != nil {
						err = multierr.Append(err, e)
					}
				} else {
					err = multierr.Append(err, e)
				}
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	// все проверки ниже только с активной SIM картой
	if b.simStatus.Load() != 1 {
		return err
	}

	// signal
	if response, e := b.client.Device.GetDeviceSignal(device.NewGetDeviceSignalParamsWithContext(ctx)); e == nil {
		const dBmPostfix = "dBm"

		if raw := strings.TrimRight(response.Payload.RSSI, dBmPostfix); raw != "" {
			if val, e := strconv.ParseInt(raw, 10, 64); e == nil {
				metricSignalRSSI.With("serial_number", sn).Set(float64(val))

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicSignalRSSI.Format(sn), val); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		if raw := strings.TrimRight(response.Payload.RSRP, dBmPostfix); raw != "" {
			if val, e := strconv.ParseInt(raw, 10, 64); e == nil {
				metricSignalRSRP.With("serial_number", sn).Set(float64(val))

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicSignalRSRP.Format(sn), val); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		if raw := strings.TrimRight(response.Payload.RSRQ, dBmPostfix); raw != "" {
			if val, e := strconv.ParseInt(raw, 10, 64); e == nil {
				metricSignalRSRQ.With("serial_number", sn).Set(float64(val))

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicSignalRSRQ.Format(sn), val); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}

		if raw := strings.TrimRight(response.Payload.SINR, dBmPostfix); raw != "" {
			if val, e := strconv.ParseInt(raw, 10, 64); e == nil {
				metricSignalSINR.With("serial_number", sn).Set(float64(val))

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicSignalSINR.Format(sn), val); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	// traffic
	if response, e := b.client.Monitoring.GetMonitoringTrafficStatistics(monitoring.NewGetMonitoringTrafficStatisticsParamsWithContext(ctx)); e == nil {
		metricTotalConnectTime.With("serial_number", sn).Set(float64(response.Payload.TotalConnectTime))
		metricTotalDownload.With("serial_number", sn).Set(float64(response.Payload.TotalDownload))
		metricTotalUpload.With("serial_number", sn).Set(float64(response.Payload.TotalUpload))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicConnectionTime.Format(sn), response.Payload.TotalConnectTime); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicConnectionDownload.Format(sn), response.Payload.TotalDownload); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicConnectionUpload.Format(sn), response.Payload.TotalUpload); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) taskCleaner(ctx context.Context) (err error) {
	if b.simStatus.Load() != 1 {
		return nil
	}

	var page int64 = 1

	for {
		paramsList := sms.NewGetSMSListParamsWithContext(ctx).
			WithRequest(&models.SMSListRequest{
				PageIndex: page,
				ReadCount: 20,
				BoxType:   1,
			})

		response, err := b.client.Sms.GetSMSList(paramsList)
		if err != nil {
			return err
		}

		if len(response.Payload.Messages) == 0 {
			return nil
		}

		for _, s := range response.Payload.Messages {
			remove := b.config.CleanerSpecial && b.checkSpecialSMS(ctx, s)
			if !remove {
				d, e := time.Parse("2006-01-02 15:04:05", s.Date)

				if e != nil {
					continue
				}

				remove = time.Since(d) > b.config.CleanerDuration
			}

			if remove {
				params := sms.NewRemoveSMSParamsWithContext(ctx)
				params.Request.Index = s.Index

				if _, err := b.client.Sms.RemoveSMS(params); err != nil {
					return err
				}

				b.Logger().Warn("Clean sms",
					"date", s.Date,
					"content", s.Content,
					"phone", s.Phone,
				)

				time.Sleep(time.Second)
			}
		}

		page++
	}
}
