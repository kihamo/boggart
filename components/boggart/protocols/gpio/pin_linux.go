// +build linux,arm

package gpio

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	g "github.com/davecheney/gpio"
)

type Pin struct {
	mutex sync.RWMutex

	pin            g.Pin
	number         int64
	status         int64
	changedAt      unsafe.Pointer
	callbackChange func(bool, time.Time, *time.Time)
}

func NewPin(number int64, mode PinMode) (GPIOPin, error) {
	var m g.Mode

	switch mode {
	case PIN_IN:
		m = g.ModeInput
	case PIN_OUT:
		m = g.ModeOutput
	case PIN_PWM:
		m = g.ModePWM
	}

	p, err := g.OpenPin(int(number), m)
	if err != nil {
		return nil, err
	}

	pin := &Pin{
		pin:    p,
		number: number,
	}

	if _, err = pin.updateAndReturn(); err != nil {
		return nil, err
	}

	if err = p.BeginWatch(g.EdgeRising, pin.watch); err != nil {
		return nil, err
	}

	return pin, nil
}

func (p *Pin) updateAndReturn() (bool, error) {
	if p.pin.Get() {
		atomic.StoreInt64(&p.status, 1)
		return true, p.pin.Err()
	}

	atomic.StoreInt64(&p.status, 0)
	return false, p.pin.Err()
}

func (p *Pin) watch() {
	currentStatus, err := p.updateAndReturn()

	if err != nil {
		fmt.Printf("Get read pin failed with error %s\n", err.Error())
		return
	}

	prevChangedAt := p.ChangedAt()
	now := time.Now()
	atomic.StorePointer(&p.changedAt, unsafe.Pointer(&now))

	p.mutex.RLock()
	cb := p.callbackChange
	p.mutex.RUnlock()

	if cb != nil {
		go func() {
			cb(currentStatus, now, prevChangedAt)
		}()
	}
}

func (p *Pin) Number() int64 {
	return p.number
}

func (p *Pin) ChangedAt() *time.Time {
	t := atomic.LoadPointer(&p.changedAt)
	return (*time.Time)(t)
}

func (p *Pin) SetCallbackChange(callback func(bool, time.Time, *time.Time)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.callbackChange = callback
}

func (p *Pin) Up() error {
	p.pin.Set()
	return p.pin.Err()
}

func (p *Pin) Down() error {
	p.pin.Clear()
	return p.pin.Err()
}

func (p *Pin) Mode() (PinMode, error) {
	m := p.pin.Mode()

	switch m {
	case g.ModeInput:
		return PIN_IN, p.pin.Err()

	case g.ModeOutput:
		return PIN_OUT, p.pin.Err()

	case g.ModePWM:
		return PIN_PWM, p.pin.Err()
	}

	return -1, p.pin.Err()
}

/**
 * Pin status:
 * 1 / true  - close
 * 0 / false - open
 */
func (p *Pin) Status() bool {
	return atomic.LoadInt64(&p.status) == 1
}
