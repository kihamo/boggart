package internal

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Device struct {
	bind        boggart.DeviceBind
	id          string
	t           string
	description string
	tags        []string
	config      interface{}
}

func (d *Device) Bind() boggart.DeviceBind {
	return d.bind
}

func (d *Device) Id() string {
	return d.id
}

func (d *Device) Type() string {
	return d.t
}

func (d *Device) Description() string {
	return d.description
}

func (d *Device) Tags() []string {
	return d.tags
}

func (d *Device) Config() interface{} {
	return d.config
}
