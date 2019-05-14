package vacuum

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
)

type CarpetMode struct {
	Enabled         bool   `json:"enabled"`
	CurrentIntegral uint64 `json:"current_integral"`
	CurrentHigh     uint64 `json:"current_high"`
	CurrentLow      uint64 `json:"current_low"`
	StallTime       uint64 `json:"stall_time"`
}

func (d *Device) CarpetMode(ctx context.Context) (result CarpetMode, err error) {
	type response struct {
		miio.Response

		Result []struct {
			CarpetMode

			Enabled uint64 `json:"enable"`
		} `json:"result"`
	}

	var reply response

	err = d.Client().Send(ctx, "get_carpet_mode", nil, &reply)
	if err == nil {
		r := &reply.Result[0]
		result.Enabled = r.Enabled == 1
		result.CurrentIntegral = r.CurrentIntegral
		result.CurrentHigh = r.CurrentHigh
		result.CurrentLow = r.CurrentLow
		result.StallTime = r.StallTime
	}

	return result, nil
}

func (d *Device) SetCarpetMode(ctx context.Context, enabled bool, integral, high, low, stallTime uint64) error {
	var reply miio.ResponseOK

	request := map[string]uint64{
		"enable":           0,
		"current_integral": integral,
		"current_high":     high,
		"current_low":      low,
		"stall_time":       stallTime,
	}

	if enabled {
		request["enable"] = 1
	}

	err := d.Client().Send(ctx, "set_carpet_mode", []interface{}{request}, &reply)
	if err != nil {
		return err
	}

	return nil
}
