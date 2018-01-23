package doors

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/davecheney/gpio"
)

/*
 * Pin state:
 * 1 / true  - close
 * 0 / false - open
 */

type Door struct {
	mutex sync.RWMutex

	pin            gpio.Pin
	status         int64
	changedAt      unsafe.Pointer
	callbackChange func(bool, *time.Time)
}

func NewDoor(pin int) (*Door, error) {
	p, err := gpio.OpenPin(pin, gpio.ModeInput)
	if err != nil {
		return nil, err
	}

	d := &Door{
		pin: p,
	}
	d.updateAndReturn()
	p.BeginWatch(gpio.EdgeFalling, d.watch)

	return d, nil
}

func (d *Door) IsOpen() bool {
	return atomic.LoadInt64(&d.status) == 0
}

func (d *Door) IsClose() bool {
	return !d.IsOpen()
}

func (d *Door) Destroy() {
	d.pin.EndWatch()
}

func (d *Door) ChangedAt() *time.Time {
	p := atomic.LoadPointer(&d.changedAt)
	return (*time.Time)(p)
}

func (d *Door) SetCallbackChange(callback func(bool, *time.Time)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.callbackChange = callback
}

func (d *Door) updateAndReturn() bool {
	if d.pin.Get() {
		atomic.StoreInt64(&d.status, 1)
		return true
	}

	atomic.StoreInt64(&d.status, 0)
	return false
}

func (d *Door) watch() {
	prevStatus := d.IsClose()
	prevChanged := d.ChangedAt()
	currentStatus := d.updateAndReturn()

	if currentStatus == prevStatus {
		return
	}

	now := time.Now()
	atomic.StorePointer(&d.changedAt, unsafe.Pointer(&now))

	d.mutex.RLock()
	cb := d.callbackChange
	d.mutex.RUnlock()

	if cb != nil {
		go func() {
			cb(currentStatus, prevChanged)
		}()
	}
}
