package openhab

import (
	"context"

	"github.com/kihamo/boggart/providers/openhab3"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	resp, err := b.provider.GetUUID(ctx)
	if err != nil {
		return err
	}

	if b.Meta().SerialNumber() == "" {
		response, err := openhab3.ParseGetUUIDResponse(resp)
		if err != nil {
			return err
		}

		if len(response.Body) != 0 {
			b.Meta().SetSerialNumber(string(response.Body))
		}
	}

	return nil
}
