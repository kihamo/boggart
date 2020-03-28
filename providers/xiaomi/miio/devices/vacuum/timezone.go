package vacuum

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

func (d *Device) Timezone(ctx context.Context) (*time.Location, error) {
	type response struct {
		miio.Response

		Result []string `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_timezone", nil, &reply)
	if err != nil {
		return nil, err
	}

	return time.LoadLocation(reply.Result[0])
}

func (d *Device) SetTimezone(ctx context.Context, zone fmt.Stringer) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "set_timezone", []string{zone.String()}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
