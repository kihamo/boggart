package connection

import (
	"bytes"
	"io"
	"net"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/protocols/connection/transport"
)

const (
	BufferSize = 512
)

var (
	readBufferPool sync.Pool
	lockerMapMutex sync.RWMutex
	lockerMap      = make(map[interface{}]*sync.Mutex)
)

func init() {
	readBufferPool.New = func() interface{} {
		buf := make([]byte, BufferSize)
		return &buf
	}
}

type Connection interface {
	io.ReadWriteCloser
	sync.Locker
	Invoke(request []byte) (response []byte, err error)
	Loop() (responses <-chan []byte, errors <-chan error, innerKill chan<- struct{}, outerDone <-chan struct{})
	ApplyOptions(options ...Option)
}

func IsTimeout(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}

type Wrapper struct {
	transport      transport.Transport
	transportMutex sync.RWMutex

	openLoops      []chan struct{}
	openLoopsMutex sync.RWMutex

	options options

	applyMutex sync.Mutex
	initDial   *atomic.Once
	lockMutex  *sync.Mutex
}

func New(transport transport.Transport, options ...Option) Connection {
	w := &Wrapper{
		transport: transport,
		openLoops: make([]chan struct{}, 0),
		initDial:  new(atomic.Once),
	}
	w.ApplyOptions(options...)

	return w
}

func (w *Wrapper) doDial() (transport transport.Transport, err error) {
	if w.options.onceInit {
		w.initDial.Do(func() {
			w.transportMutex.RLock()
			t := w.transport
			w.transportMutex.RUnlock()

			transport, err = t.Dial()

			if err == nil {
				w.transportMutex.Lock()
				w.transport = transport
				w.transportMutex.Unlock()
			}
		})

		if err != nil {
			w.initDial.Reset()
		}
	} else {
		w.transportMutex.RLock()
		t := w.transport
		w.transportMutex.RUnlock()

		transport, err = t.Dial()

		if err == nil {
			w.transportMutex.Lock()
			w.transport = transport
			w.transportMutex.Unlock()
		}
	}

	w.transportMutex.RLock()
	defer w.transportMutex.RUnlock()

	return w.transport, err
}

func (w *Wrapper) Read(p []byte) (n int, err error) {
	conn, err := w.doDial()
	if err != nil {
		return n, err
	}

	w.Lock()
	defer w.Unlock()

	n, err = conn.Read(p)

	if n > 0 && w.options.dumpRead != nil {
		w.options.dumpRead(p[:n])
	}

	return n, err
}

func (w *Wrapper) Write(p []byte) (n int, err error) {
	conn, err := w.doDial()
	if err != nil {
		return n, err
	}

	w.Lock()
	defer w.Unlock()

	n, err = conn.Write(p)

	if n > 0 && w.options.dumpWrite != nil {
		w.options.dumpWrite(p[:n])
	}

	return n, err
}

func (w *Wrapper) Close() error {
	w.openLoopsMutex.RLock()
	for _, ch := range w.openLoops {
		ch <- struct{}{}
	}
	w.openLoopsMutex.RUnlock()

	w.transportMutex.RLock()
	t := w.transport
	w.transportMutex.RUnlock()

	if t != nil {
		return t.Close()
	}

	return nil
}

func (w *Wrapper) Lock() {
	if w.lockMutex != nil {
		w.lockMutex.Lock()
	}
}

func (w *Wrapper) Unlock() {
	if w.lockMutex != nil {
		w.lockMutex.Unlock()
	}
}

func (w *Wrapper) Invoke(request []byte) (response []byte, err error) {
	conn, err := w.doDial()
	if err != nil {
		return nil, err
	}

	w.Lock()
	defer w.Unlock()

	var n int

	n, err = conn.Write(request)
	if err != nil {
		return nil, err
	}

	if n > 0 && w.options.dumpWrite != nil {
		w.options.dumpWrite(request[:n])
	}

	buffer := bytes.NewBuffer(nil)

	buf := readBufferPool.Get().(*[]byte)
	defer readBufferPool.Put(buf)

	for {
		n, err = conn.Read(*buf)

		if n > 0 {
			if w.options.dumpRead != nil {
				w.options.dumpRead((*buf)[:n])
			}

			buffer.Write((*buf)[:n])

			if w.options.readCheck != nil {
				if !w.options.readCheck(buffer.Bytes()) {
					break
				}
			} else if n < len(*buf) {
				// FIXME: если ответ будет ровно равен len(*buf) то следующий цикл приведет к блокировке
				break
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (w *Wrapper) Loop() (<-chan []byte, <-chan error, chan<- struct{}, <-chan struct{}) {
	response := make(chan []byte)
	errors := make(chan error)
	innerKill := make(chan struct{}, 1)
	outerDone := make(chan struct{}, 1)

	go func() {
		buf := make([]byte, BufferSize)

		for {
			select {
			case <-innerKill:
				close(response)
				close(errors)
				close(outerDone)

				w.openLoopsMutex.Lock()
				for i, ch := range w.openLoops {
					if ch == innerKill {
						w.openLoops = append(w.openLoops[:i], w.openLoops[i+1:]...)
						break
					}
				}
				w.openLoopsMutex.Unlock()

				return

			default:
				n, err := w.Read(buf)

				if n > 0 {
					response <- append([]byte(nil), buf[:n]...)
				}

				if err != nil && !IsTimeout(err) {
					errors <- err
				}
			}
		}
	}()

	w.openLoopsMutex.Lock()
	w.openLoops = append(w.openLoops, innerKill)
	w.openLoopsMutex.Unlock()

	return response, errors, innerKill, outerDone
}

func (w *Wrapper) ApplyOptions(options ...Option) {
	for _, option := range options {
		option.apply(&w.options)
	}

	var lockMutex *sync.Mutex

	if w.options.lockGlobal {
		addr, ok := w.transport.Options()["address"]
		if ok {
			lockerMapMutex.Lock()
			if _, ok := lockerMap[addr]; !ok {
				lockerMap[addr] = new(sync.Mutex)
			}
			lockMutex = lockerMap[addr]
			lockerMapMutex.Unlock()
		} else {
			lockMutex = new(sync.Mutex)
		}
	} else if w.options.lockLocal {
		lockMutex = new(sync.Mutex)
	}

	w.applyMutex.Lock()
	defer w.applyMutex.Unlock()

	if lockMutex != w.lockMutex {
		w.lockMutex = lockMutex
	}
}
