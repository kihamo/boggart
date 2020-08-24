package atomic

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) != 0 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()

	atomic.StoreUint32(&o.done, 0)
}

func (o *Once) IsDone() bool {
	return atomic.LoadUint32(&o.done) != 0
}
