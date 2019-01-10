package internal

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Device struct {
	bind        boggart.DeviceBind
	id          string
	description string
	tags        []string
	config      map[string]interface{}
}

func (d *Device) Bind() boggart.DeviceBind {
	return d.bind
}

func (d *Device) Id() string {
	return d.id
}

func (d *Device) Description() string {
	return d.description
}

func (d *Device) Tags() []string {
	return d.tags
}

func (d *Device) Config() map[string]interface{} {
	return d.config
}
