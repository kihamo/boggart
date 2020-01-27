package hikvision

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/hikvision/client/system"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	deviceInfo, err := b.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)

	if err == nil {
		if deviceInfo.Payload.SerialNumber == "" {
			return errors.New("device returns empty serial number")
		}

		if deviceInfo.Payload.MacAddress == "" {
			return errors.New("device returns empty MAC address")
		}

		b.Meta().SetSerialNumber(deviceInfo.Payload.SerialNumber)

		err = b.Meta().SetMACAsString(deviceInfo.Payload.MacAddress)
	}

	return err
}
