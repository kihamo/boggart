package connection

import (
	"io"
	"sync"
)

type ObserveFunc struct {
	f func([]byte, error)
}

type Observer interface {
	Packet(response []byte)
	Error(err error)
}

type ObserverConnection interface {
	Invoker

	Attach(Observer)
	Detach(Observer)
}

type observerConnection struct {
	observers     []Observer
	observersLock sync.RWMutex

	conn     Looper
	connOnce sync.Once
}

func (o *ObserveFunc) Packet(response []byte) {
	o.f(response, nil)
}

func (o *ObserveFunc) Error(err error) {
	o.f(nil, err)
}

func ObserverFunc(f func([]byte, error)) Observer {
	return &ObserveFunc{f: f}
}

func NewObserverConnection(conn Conn) ObserverConnection {
	if i, ok := conn.(ObserverConnection); ok {
		return i
	}

	return &observerConnection{
		conn: NewLooper(conn),
	}
}

func (m *observerConnection) initLoop() {
	m.connOnce.Do(func() {
		go m.loop()
	})
}

func (m *observerConnection) loop() {
	chReceivePackets, chReceiveErrors, chConnKill, chConnDone := m.conn.Loop()

	closer := func() {
		chConnKill <- struct{}{} // kill connect
	}

	for {
		select {
		case packet := <-chReceivePackets:
			go func(data []byte) {
				m.observersLock.RLock()
				defer m.observersLock.RUnlock()

				for _, o := range m.observers {
					o.Packet(append([]byte(nil), data...))
				}
			}(packet)

		case err := <-chReceiveErrors:
			if err == nil || err == io.EOF {
				continue
			}

			go func(e error) {
				m.observersLock.RLock()
				defer m.observersLock.RUnlock()

				for _, o := range m.observers {
					o.Error(e)
				}
			}(err)

			closer()

		case <-chConnDone:
			return
		}
	}
}

func (m *observerConnection) Attach(o Observer) {
	m.observersLock.Lock()
	m.observers = append(m.observers, o)
	m.observersLock.Unlock()

	m.initLoop()
}

func (m *observerConnection) Detach(o Observer) {
	m.observersLock.Lock()
	defer m.observersLock.Unlock()

	for i := len(m.observers) - 1; i >= 0; i-- {
		if m.observers[i] == o {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			// watcher.close()
		}
	}
}

func (m *observerConnection) Invoke(request []byte) (response []byte, err error) {
	if locker, ok := m.conn.(sync.Locker); ok {
		locker.Lock()
		defer locker.Unlock()
	}

	done := make(chan struct{}, 1)

	observer := ObserverFunc(func(packet []byte, e error) {
		response = packet
		err = e

		done <- struct{}{}

	})
	m.Attach(observer)
	defer m.Detach(observer)

	if _, err := m.conn.Write(request); err != nil {
		return response, err
	}

	<-done

	return response, err
}

func (m *observerConnection) Read(p []byte) (n int, err error) {
	m.initLoop()

	return m.conn.Read(p)
}

func (m *observerConnection) Write(p []byte) (n int, err error) {
	m.initLoop()

	return m.conn.Write(p)
}

func (m *observerConnection) Close() error {
	return m.conn.Close()
}
