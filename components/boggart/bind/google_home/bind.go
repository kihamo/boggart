package google_home

import (
	"context"
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/google/home"
	"github.com/kihamo/boggart/providers/google/home/client/info"
)

type Bind struct {
	di.MetaBind
	di.LoggerBind
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

	if mac, err := net.ParseMAC(response.Payload.MacAddress); err == nil {
		b.Meta().SetMAC(mac)
	}

	return response.Payload.MacAddress, nil
}
