package home

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/google/home"
	"github.com/kihamo/boggart/providers/google/home/client/info"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.ProbesBind

	provider *home.Client
}

func (b *Bind) GetMAC(ctx context.Context) (string, error) {
	response, err := b.provider.Info.GetEurekaInfo(info.NewGetEurekaInfoParamsWithContext(ctx).
		WithOptions(home.EurekaInfoOptionDetail.Value()).
		WithParams(home.EurekaInfoParamDeviceInfo.Value()))

	if err != nil {
		return "", err
	}

	if response.Payload == nil || response.Payload.MacAddress == "" {
		return "", errors.New("MAC address not found")
	}

	if err := b.Meta().SetMACAsString(response.Payload.MacAddress); err != nil {
		return "", err
	}

	return response.Payload.MacAddress, nil
}
