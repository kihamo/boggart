package bluetooth

import (
	"context"
	"errors"

	"github.com/go-ble/ble"
)

type Dummy struct {
}

func NewDummyDevice(_ ...ble.Option) (ble.Device, error) {
	return &Dummy{}, nil
}

func (d *Dummy) AddService(svc *ble.Service) error {
	return nil
}

func (d *Dummy) RemoveAllServices() error {
	return nil
}

func (d *Dummy) SetServices(svcs []*ble.Service) error {
	return nil
}

func (d *Dummy) Stop() error {
	return nil
}

func (d *Dummy) Advertise(_ context.Context, _ ble.Advertisement) error {
	return nil
}

func (d *Dummy) AdvertiseNameAndServices(_ context.Context, _ string, _ ...ble.UUID) error {
	return nil
}

func (d *Dummy) AdvertiseMfgData(_ context.Context, _ uint16, _ []byte) error {
	return nil
}

func (d *Dummy) AdvertiseServiceData16(_ context.Context, _ uint16, _ []byte) error {
	return nil
}

func (d *Dummy) AdvertiseIBeaconData(_ context.Context, _ []byte) error {
	return nil
}

func (d *Dummy) AdvertiseIBeacon(_ context.Context, _ ble.UUID, _, _ uint16, _ int8) error {
	return nil
}

func (d *Dummy) Scan(_ context.Context, _ bool, _ ble.AdvHandler) error {
	return nil
}

func (d *Dummy) Dial(_ context.Context, _ ble.Addr) (ble.Client, error) {
	return nil, errors.New("not implemented")
}
