package hilink

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/hilink"
	"github.com/kihamo/boggart/providers/hilink/client/device"
	"github.com/kihamo/boggart/providers/hilink/client/ussd"
	"github.com/kihamo/boggart/providers/hilink/models"
)

var (
	cmdRegexp = regexp.MustCompile(`^(exec)?\s*(?P<command>.*)$`)
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config                    *Config
	client                    *hilink.Client
	operator                  *atomic.String
	limitInternetTrafficIndex *atomic.Int64
	simStatus                 *atomic.Uint32
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
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	op, err := b.operatorSettings()
	if err != nil {
		return -1, err
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

func (b *Bind) operatorSettings() (*operator, error) {
	label := b.operator.Load()

	switch strings.ToLower(label) {
	case "tele2", "tele2 ru":
		return operatorTele2, nil
	}

	return nil, errors.New("operator " + label + " settings isn't found")
}

func (b *Bind) checkSpecialSMS(ctx context.Context, sms *models.SMSListMessagesItems0) bool {
	op, err := b.operatorSettings()
	if err != nil {
		return false
	}

	sn := b.Meta().SerialNumber()

	// limit traffic
	match := op.SMSLimitTrafficRegexp.FindStringSubmatch(sms.Content)
	if len(match) > 0 {
		var (
			value float64
			found bool
		)

		for i, name := range op.SMSLimitTrafficRegexp.SubexpNames() {
			if !strings.HasPrefix(name, "value") {
				continue
			}

			found = true

			if v, err := strconv.ParseFloat(match[i], 64); err == nil {
				value += v
			}
		}

		if found {
			if sms.Index > b.limitInternetTrafficIndex.Load() {
				value *= op.SMSLimitTrafficFactor

				metricLimitInternetTraffic.With("serial_number", sn).Set(value)
				b.MQTT().PublishAsync(ctx, b.config.TopicLimitInternetTraffic.Format(sn), uint64(value))

				b.limitInternetTrafficIndex.Set(sms.Index)
			}

			return true
		}
	}

	// special commands
	if b.config.SMSCommandsEnabled {
		match = cmdRegexp.FindStringSubmatch(sms.Content)
		if len(match) > 0 {
			for i, name := range cmdRegexp.SubexpNames() {
				if name == "command" {
					cmd := strings.ToLower(match[i])

					switch cmd {
					case "reboot":
						// поддерживаемые команды
					default:
						return false
					}

					// нет разрешенных телефонов с которых можно принимать команды, но команда найдена
					if len(b.config.SMSCommandsAllowedPhones) == 0 {
						return true
					}

					var allowed bool

					for _, phone := range b.config.SMSCommandsAllowedPhones {
						if strings.Compare(phone, sms.Phone) == 0 {
							allowed = true
							break
						}
					}

					// выполнение команд с номера не разрешено
					if !allowed {
						b.Logger().Warnf("Unauthorized execute command %s from phone number %s ", cmd, sms.Phone)
						return true
					}

					if cmd == "reboot" {
						params := device.NewDeviceControlParams()
						params.Request.Control = 1

						if _, err := b.client.Device.DeviceControl(params); err != nil {
							b.Logger().Errorf("Reboot failed with error %s", err.Error())
						}
					}

					return true
				}
			}
		}
	}

	return false
}
