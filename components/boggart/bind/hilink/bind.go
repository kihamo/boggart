package hilink

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/providers/hilink"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/ussd"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/models"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config                    *Config
	client                    *hilink.Client
	operator                  *atomic.String
	limitInternetTrafficIndex *atomic.Int64
}

func (b *Bind) USSD(ctx context.Context, content string) (string, error) {
	if content == "" {
		return "", nil
	}

	params := ussd.NewSendUSSDParamsWithContext(ctx).
		WithRequest(&models.USSD{
			Content: content,
		})

	_, err := b.client.Ussd.SendUSSD(params)
	if err != nil {
		return "", err
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		default:
			response, err := b.client.Ussd.GetUSSD(ussd.NewGetUSSDParamsWithContext(ctx))
			if err == nil && response.Payload.Content != "" {
				return response.Payload.Content, nil
			}

			time.Sleep(time.Second)
		}
	}

	return "", err
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	op := b.operatorSettings()
	if op == nil {
		return -1, errors.New("operator settings isn't found")
	}

	content, err := b.USSD(ctx, op.BalanceUSSD)
	if err != nil {
		return -1, err
	}

	match := op.BalanceRegexp.FindStringSubmatch(content)
	for i, name := range op.BalanceRegexp.SubexpNames() {
		if name == "value" {
			return strconv.ParseFloat(match[i], 64)
		}
	}

	return 0, errors.New("balance not found")
}

func (b *Bind) operatorSettings() *operator {
	switch strings.ToLower(b.operator.Load()) {
	case "tele2", "tele2 ru":
		return operatorTele2
	}

	return nil
}

func (b *Bind) checkSpecialSMS(ctx context.Context, sms *models.SMSListMessagesItems0) (result bool) {
	op := b.operatorSettings()
	if op == nil {
		return result
	}

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)

	match := op.SMSLimitTrafficRegexp.FindStringSubmatch(sms.Content)
	for i, name := range op.SMSLimitTrafficRegexp.SubexpNames() {
		if name == "value" {
			result = true

			if sms.Index > b.limitInternetTrafficIndex.Load() {
				if value, err := strconv.ParseFloat(match[i], 64); err == nil {
					value *= op.SMSLimitTrafficFactor

					metricLimitInternetTraffic.With("serial_number", sn).Set(value)
					b.MQTTPublishAsync(ctx, MQTTPublishTopicLimitInternetTraffic.Format(snMQTT), uint64(value))

					b.limitInternetTrafficIndex.Set(sms.Index)
				}
			}

			break
		}
	}

	return result
}
