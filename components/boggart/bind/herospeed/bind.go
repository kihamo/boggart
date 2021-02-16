package herospeed

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/herospeed"
	"go.uber.org/multierr"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	config *Config
	client *herospeed.Client
}

func (b *Bind) GetSerialNumber(ctx context.Context) (string, error) {
	configuration, err := b.client.Configuration(ctx)
	if err != nil {
		return "", err
	}

	sn, ok := configuration["serialnumber"]
	if !ok || len(sn) == 0 {
		return "", errors.New("device returns empty serial number")
	}

	if sn != b.Meta().SerialNumber() {
		b.Meta().SetSerialNumber(sn)

		if mac, ok := configuration["macip"]; ok {
			if err := b.Meta().SetMACAsString(mac); err != nil {
				return "", err
			}
		}

		var mqttError error

		if model, ok := configuration["modelname"]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(sn), model); e != nil {
				mqttError = multierr.Append(mqttError, e)
			}
		}

		if fw, ok := configuration["firmwareversion"]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(sn), fw); e != nil {
				mqttError = multierr.Append(mqttError, e)
			}
		}

		if mqttError != nil {
			b.Logger().Error(mqttError.Error())
		}
	}

	return sn, nil
}
