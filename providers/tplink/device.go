package tplink

import (
	"context"

	"github.com/kihamo/boggart/providers/tplink/internal"
)

type Kind int

const (
	KindBuld Kind = 1
)

type Device interface {
	GetDeviceInfo(context.Context) (interface{}, error)
}

func NewDevice(addr, email, password string) Device {
	return internal.NewDevice(addr, email, password)
}
