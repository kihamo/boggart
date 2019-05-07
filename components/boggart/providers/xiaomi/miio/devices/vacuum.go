package devices

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
)

// https://github.com/marcelrv/XiaomiRobotVacuumProtocol

type Vacuum struct {
	miio.Device
}

func NewVacuum(address, token string) *Vacuum {
	d := &Vacuum{
		Device: *miio.NewDevice(address, token),
	}

	return d
}

func (d *Vacuum) SerialNumber() (string, error) {
	type response struct {
		miio.Response

		Result []struct {
			SerialNumber string `json:"serial_number"`
		} `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_serial_number", nil, &reply)
	if err != nil {
		return "", err
	}

	return reply.Result[0].SerialNumber, nil
}

func (d *Vacuum) SoundVolume() (uint64, error) {
	type response struct {
		miio.Response

		Result []uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_sound_volume", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SetSoundVolume(volume uint64) error {
	if volume > 100 {
		volume = 100
	}

	var reply miio.ResponseOK

	err := d.Client().Send("change_sound_volume", []uint64{volume}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
