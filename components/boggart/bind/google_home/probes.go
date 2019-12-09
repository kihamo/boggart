package google_home

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/google/home"
	"github.com/kihamo/boggart/providers/google/home/client/info"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	response, err := b.ClientGoogleHome().Info.GetEurekaInfo(info.NewGetEurekaInfoParams().
		WithOptions(home.EurekaInfoOptionDetail.Value()).
		WithParams(home.EurekaInfoParamDeviceInfo.Value()))

	if err != nil {
		return err
	}

	if b.SerialNumber() == "" {
		if response.Payload == nil || response.Payload.MacAddress == "" {
			return errors.New("MAC address not found")
		}

		b.SetSerialNumber(response.Payload.MacAddress)
	}

	return nil
}
