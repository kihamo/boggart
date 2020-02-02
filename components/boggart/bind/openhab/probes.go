package openhab

import (
	"context"

	"github.com/kihamo/boggart/providers/openhab/client/uuid"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	response, err := b.provider.UUID.GetInstanceUUID(uuid.NewGetInstanceUUIDParamsWithContext(ctx))
	if err == nil {
		b.Meta().SetSerialNumber(response.Payload)
	}

	return err
}
