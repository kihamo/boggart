package vacuum

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

type DoNotDisturb struct {
	Enabled     bool   `json:"enabled"`
	StartHour   uint64 `json:"start_hour"`
	StartMinute uint64 `json:"start_minute"`
	EndHour     uint64 `json:"end_hour"`
	EndMinute   uint64 `json:"end_minute"`
}

func (d *Device) DoNotDisturb(ctx context.Context) (result DoNotDisturb, err error) {
	type response struct {
		miio.Response

		Result []struct {
			DoNotDisturb

			Enabled uint64 `json:"enabled"`
		} `json:"result"`
	}

	var reply response

	err = d.Client().Send(ctx, "get_dnd_timer", nil, &reply)
	if err == nil {
		r := &reply.Result[0]
		result.Enabled = r.Enabled == 1
		result.StartHour = r.StartHour
		result.StartMinute = r.StartMinute
		result.EndHour = r.EndHour
		result.EndMinute = r.EndMinute
	}

	return result, err
}

func (d *Device) DisableDoNotDisturb(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "close_dnd_timer", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) SetDoNotDisturb(ctx context.Context, startHour, startMinute, endHour, endMinute uint64) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "set_dnd_timer", []uint64{startHour, startMinute, endHour, endMinute}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
