package samsung_tizen

import (
	"context"
	"net"
	"strings"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/samsung/tv"
	"go.uber.org/multierr"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config *Config
	client *tv.ApiV2
}

func (b *Bind) GetSerialNumber(ctx context.Context) (sn string, err error) {
	info, err := b.client.Device(ctx)
	if err != nil {
		return "", err
	}

	parts := strings.Split(info.ID, ":")
	if len(parts) > 1 {
		sn = parts[1]
	} else {
		sn = info.ID
	}

	if sn != b.Meta().SerialNumber() {
		b.Meta().SetSerialNumber(sn)
		var mqttError error

		if mac, e := net.ParseMAC(info.Device.WifiMac); e == nil {
			b.Meta().SetMAC(mac)
		} else {
			mqttError = multierr.Append(mqttError, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDeviceID.Format(sn), info.Device.ID); e != nil {
			mqttError = multierr.Append(mqttError, e)
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicDeviceModelName.Format(sn), info.Device.Name); e != nil {
			mqttError = multierr.Append(mqttError, e)
		}

		if mqttError != nil {
			b.Logger().Error(mqttError.Error())
		}
	}

	return sn, nil
}
