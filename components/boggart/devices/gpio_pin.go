package devices

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/stianeikeland/go-rpio"
)

var initGPIO sync.Once

type GPIOPin struct {
	boggart.DeviceBase

	pin   rpio.Pin
	value uint64
}

func NewGPIOPin(number uint64, mode rpio.Mode) *GPIOPin {
	initGPIO.Do(func() {
		if err := rpio.Open(); err != nil {
			fmt.Println("Failed init GPIO with error", err.Error())
		}
	})

	device := &GPIOPin{
		pin: rpio.Pin(number),
	}
	device.pin.Mode(mode)

	modeAsString := "unknown"
	switch mode {
	case rpio.Input:
		modeAsString = "input"
	case rpio.Output:
		modeAsString = "output"
	case rpio.Clock:
		modeAsString = "clock"
	case rpio.Pwm:
		modeAsString = "pwm"
	}

	device.Init()
	device.SetDescription(fmt.Sprintf("GPIO pin #%d with %s mode", number, modeAsString))

	if mode == rpio.Input {
		go func() {
			device.edgeDetected()
		}()
	}

	return device
}

func (d *GPIOPin) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeGPIOPin,
	}
}

func (d *GPIOPin) Ping(_ context.Context) bool {
	return true
}

func (d *GPIOPin) High() {
	d.pin.High()
}

func (d *GPIOPin) Low() {
	d.pin.Low()
}

func (d *GPIOPin) edgeDetected() {
	d.pin.Detect(rpio.AnyEdge)

	for {
		prev := atomic.LoadUint64(&d.value)
		current := uint64(d.pin.Read())

		if prev != current {
			atomic.StoreUint64(&d.value, current)
			d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, uint64(d.pin), current)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
