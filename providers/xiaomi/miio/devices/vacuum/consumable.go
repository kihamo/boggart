package vacuum

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

const (
	ConsumableFilter    consumable = "filter_work_time"
	ConsumableBrushMain consumable = "main_brush_work_time"
	ConsumableBrushSide consumable = "side_brush_work_time"
	ConsumableSensor    consumable = "sensor_dirty_time"

	ConsumableLifetimeFilter    time.Duration = 150
	ConsumableLifetimeBrushMain time.Duration = 300
	ConsumableLifetimeBrushSide time.Duration = 200
	ConsumableLifetimeSensor    time.Duration = 30
)

type consumable string

// nolint:golint
func (d *Device) Consumables(ctx context.Context) (map[consumable]time.Duration, error) {
	var response struct {
		miio.Response

		Result []map[string]uint64 `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "get_consumable", nil, &response)
	if err != nil {
		return nil, err
	}

	consumables := make(map[consumable]time.Duration, len(response.Result[0]))
	for n, v := range response.Result[0] {
		consumables[consumable(n)] = time.Duration(v) * time.Second
	}

	return consumables, nil
}

func (d *Device) ConsumableReset(ctx context.Context, consumable consumable) error {
	var response miio.ResponseOK

	err := d.Client().CallRPC(ctx, "reset_consumable", []string{string(consumable)}, &response)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(response) {
		return errors.New("device return not OK response")
	}

	return nil
}
