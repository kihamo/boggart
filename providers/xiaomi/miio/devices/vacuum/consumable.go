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
	type response struct {
		miio.Response

		Result []map[string]uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_consumable", nil, &reply)
	if err != nil {
		return nil, err
	}

	consumables := make(map[consumable]time.Duration, len(reply.Result[0]))
	for n, v := range reply.Result[0] {
		consumables[consumable(n)] = time.Duration(v) * time.Second
	}

	return consumables, nil
}

func (d *Device) ConsumableReset(ctx context.Context, consumable consumable) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "reset_consumable", []string{string(consumable)}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
