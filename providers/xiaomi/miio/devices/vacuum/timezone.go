package vacuum

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

func (d *Device) Timezone(ctx context.Context) (*time.Location, error) {
	var response struct {
		miio.Response

		Result []string `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "get_timezone", nil, &response)
	if err != nil {
		return nil, err
	}

	return time.LoadLocation(response.Result[0])
}

func (d *Device) SetTimezone(ctx context.Context, zone fmt.Stringer) error {
	var response miio.ResponseOK

	err := d.Client().CallRPC(ctx, "set_timezone", []string{zone.String()}, &response)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(response) {
		return errors.New("device return not OK response")
	}

	return nil
}
