package doors

import (
	"sync"
	"sync/atomic"

	"github.com/davecheney/gpio"
)

const (
	OPEN  = false
	CLOSE = true
)

type Door struct {
	mutex sync.RWMutex

	pin      gpio.Pin
	status   int64
	callback func(bool)
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

func (d *Door) SetCallback(callback func(bool)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.callback = callback
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
	val := d.updateAndReturn()

	d.mutex.RLock()
	cb := d.callback
	d.mutex.RUnlock()

	if cb != nil {
		go func() {
			cb(val)
		}()
	}
}
