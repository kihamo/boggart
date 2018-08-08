package devices

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/pin"
)

type GPIOMode int64

const (
	GPIOModeDefault GPIOMode = iota
	GPIOModeIn
	GPIOModeOut
)

type GPIOPin struct {
	value uint64

	boggart.DeviceBase
	pin  pin.Pin
	mode GPIOMode
}

func NewGPIOPin(p pin.Pin, m GPIOMode) *GPIOPin {
	device := &GPIOPin{
		pin:  p,
		mode: m,
	}

	device.Init()
	device.SetDescription(fmt.Sprintf("%s %s", p.Name(), p.Function()))

	if _, ok := p.(gpio.PinIn); ok {
		go func() {
			device.waitForEdge()
		}()
	}

	return device
}

func (d *GPIOPin) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeGPIO,
	}
}

func (d *GPIOPin) Ping(_ context.Context) bool {
	return true
}

func (d *GPIOPin) Mode() GPIOMode {
	return d.mode
}

func (d *GPIOPin) High() error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		return g.Out(gpio.High)
	}

	return nil
}

func (d *GPIOPin) Low() error {
	if d.Mode() == GPIOModeIn {
		return nil
	}

	if g, ok := d.pin.(gpio.PinOut); ok {
		return g.Out(gpio.Low)
	}

	return nil
}

func (d *GPIOPin) Read() bool {
	if d.Mode() == GPIOModeOut {
		return false
	}

	if g, ok := d.pin.(gpio.PinIn); ok {
		return g.Read() == gpio.High
	}

	return false
}

func (d *GPIOPin) waitForEdge() {
	p := d.pin.(gpio.PinIn)
	p.In(gpio.PullNoChange, gpio.BothEdges)

	for p.WaitForEdge(-1) {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.pin.Number(), d.Read())
	}
}
