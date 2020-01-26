package herospeed

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/herospeed"
	"go.uber.org/multierr"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

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

		if model, ok := configuration["modelname"]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateModel.Format(sn), model); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if fw, ok := configuration["firmwareversion"]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateFirmwareVersion.Format(sn), fw); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return sn, err
}
