package vacuum

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

const (
	FanPowerQuiet    uint64 = 38
	FanPowerBalanced uint64 = 60
	FanPowerTurbo    uint64 = 75
	FanPowerMax      uint64 = 100
	FanPowerMob      uint64 = 105
)

func (d *Device) FanPower(ctx context.Context) (uint32, error) {
	var response struct {
		miio.Response

		Result []uint32 `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "get_custom_mode", nil, &response)
	if err != nil {
		return 0, err
	}

	return response.Result[0], nil
}

func (d *Device) SetFanPower(ctx context.Context, power uint64) error {
	if power > 105 {
		power = 105
	} else if power < 1 {
		power = 1
	}

	var response miio.ResponseOK

	err := d.Client().CallRPC(ctx, "set_custom_mode", []uint64{power}, &response)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(response) {
		return errors.New("device return not OK response")
	}

	return nil
}
