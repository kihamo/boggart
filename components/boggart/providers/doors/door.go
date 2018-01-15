package doors

import (
	"sync"

	"github.com/davecheney/gpio"
)

const (
	OPEN  = false
	CLOSE = true
)

type Door struct {
	mutex sync.RWMutex

	pin      gpio.Pin
	// FIXME: sync/atomic
	status   bool
	callback func(bool)
}

func NewDoor(pin int) (*Door, error) {
	p, err := gpio.OpenPin(pin, gpio.ModeInput)
	if err != nil {
		return nil, err
	}

	d := &Door{
		pin:    p,
		status: p.Get(),
	}
	p.BeginWatch(gpio.EdgeFalling, d.watch)

	return d, nil
}

func (d *Door) IsOpen() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.status == OPEN
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

func (d *Door) watch() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.status = d.pin.Get()

	if d.callback != nil {
		go func() {
			d.callback(d.status)
		}()
	}
}
