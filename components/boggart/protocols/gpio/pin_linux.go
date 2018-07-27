// +build linux,arm

package gpio

import (
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
	callbackChange func(bool, *time.Time)
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
	pin.updateAndReturn()
	p.BeginWatch(g.EdgeRising, pin.watch)

	return pin, nil
}

func (p *Pin) updateAndReturn() bool {
	if p.pin.Get() {
		atomic.StoreInt64(&p.status, 1)
		return true
	}

	atomic.StoreInt64(&p.status, 0)
	return false
}

func (p *Pin) watch() {
	prevStatus := p.Status()
	prevChanged := p.ChangedAt()
	currentStatus := p.updateAndReturn()

	if currentStatus == prevStatus {
		return
	}

	now := time.Now()
	atomic.StorePointer(&p.changedAt, unsafe.Pointer(&now))

	p.mutex.RLock()
	cb := p.callbackChange
	p.mutex.RUnlock()

	if cb != nil {
		go func() {
			cb(currentStatus, prevChanged)
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

func (p *Pin) SetCallbackChange(callback func(bool, *time.Time)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.callbackChange = callback
}

/**
 * Pin status:
 * 1 / true  - close
 * 0 / false - open
 */
func (p *Pin) Status() bool {
	return atomic.LoadInt64(&p.status) == 1
}
