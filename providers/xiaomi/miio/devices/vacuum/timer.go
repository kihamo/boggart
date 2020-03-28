package vacuum

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
)

type Timer struct {
	ID       uint64
	Enabled  bool
	Cron     string
	Action   string
	FanPower uint64
}

func (d *Device) Timers(ctx context.Context) ([]Timer, error) {
	type response struct {
		miio.Response

		Result [][]interface{} `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_timer", nil, &reply)
	if err != nil {
		return nil, err
	}

	result := make([]Timer, 0, len(reply.Result))

	for _, row := range reply.Result {
		t := Timer{}

		for i, v := range row {
			switch i {
			case 0:
				t.ID, _ = strconv.ParseUint(v.(string), 10, 64)

			case 1:
				t.Enabled = strings.EqualFold(v.(string), "on")

			case 2:
				values := v.([]interface{})
				t.Cron = values[0].(string)

				command := values[1].([]interface{})
				t.Action = command[0].(string)
				t.FanPower = uint64(command[1].(float64))
			}
		}

		result = append(result, t)
	}

	return result, nil
}

func (d *Device) DisableTimer(ctx context.Context, timerID uint64, value bool) error {
	val := "off"
	if value {
		val = "on"
	}

	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "upd_timer", []string{strconv.FormatUint(timerID, 10), val}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) RemoveTimer(ctx context.Context, timerID uint64) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "del_timer", []string{strconv.FormatUint(timerID, 10)}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) SetTimer(ctx context.Context, cron, action string, fanPower uint64) error {
	var reply miio.ResponseOK

	request := [][]interface{}{{
		strconv.FormatInt(time.Now().Unix(), 10),
		[]interface{}{
			cron,
			[]interface{}{action, fanPower},
		},
	}}

	err := d.Client().Send(ctx, "set_timer", request, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
