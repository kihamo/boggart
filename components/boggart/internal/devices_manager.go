package internal

import (
	"context"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

type DevicesManager struct {
	mutex          sync.RWMutex
	storage        *sync.Map
	workers        workers.Component
	logger         logger.Logger
	listeners      *manager.ListenersManager
	chanChecker    chan struct{}
	tickerChecker  *w.Ticker
	timeoutChecker time.Duration
}

func NewDevicesManager(workers workers.Component) *DevicesManager {
	m := &DevicesManager{
		storage:       new(sync.Map),
		workers:       workers,
		logger:        logger.NopLogger,
		listeners:     manager.NewListenersManager(),
		chanChecker:   make(chan struct{}, 1),
		tickerChecker: w.NewTicker(time.Minute),
		// TODO: setter for this option
		timeoutChecker: time.Second * 2,
	}

	go m.doCheck()

	return m
}

func (m *DevicesManager) SetLogger(logger logger.Logger) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.logger = logger
}

func (m *DevicesManager) Logger() logger.Logger {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.logger
}

func (m *DevicesManager) Register(device boggart.Device) string {
	id := uuid.New()
	m.RegisterWithID(id, device)
	return id
}

func (m *DevicesManager) RegisterWithID(id string, device boggart.Device) {
	m.storage.Store(id, device)
	m.listeners.AsyncTrigger(boggart.DeviceEventDeviceRegister, device, id)

	tasks := device.Tasks()
	if len(tasks) > 0 {
		for _, task := range tasks {
			m.workers.AddTask(task)
		}
	}
}

func (m *DevicesManager) Device(id string) boggart.Device {
	if d, ok := m.storage.Load(id); ok {
		return d.(boggart.Device)
	}

	return nil
}

func (m *DevicesManager) Devices() map[string]boggart.Device {
	devices := make(map[string]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		devices[key.(string)] = device.(boggart.Device)
		return true
	})

	return devices
}

func (m *DevicesManager) DevicesByTypes(types []boggart.DeviceType) map[string]boggart.Device {
	if len(types) == 0 {
		return m.Devices()
	}

	devices := make(map[string]boggart.Device, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		d := device.(boggart.Device)

		for _, t1 := range d.Types() {
			for _, t2 := range types {
				if t1.String() == t2.String() {
					devices[key.(string)] = d
					return true
				}
			}
		}

		return true
	})

	return devices
}

func (m *DevicesManager) Describe(ch chan<- *snitch.Description) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if !device.(boggart.Device).IsEnabled() {
			return true
		}

		if collector, ok := device.(snitch.Collector); ok {
			collector.Describe(ch)
		}

		return true
	})
}

func (m *DevicesManager) Collect(ch chan<- snitch.Metric) {
	m.storage.Range(func(_ interface{}, device interface{}) bool {
		if !device.(boggart.Device).IsEnabled() {
			return true
		}

		if collector, ok := device.(snitch.Collector); ok {
			collector.Collect(ch)
		}

		return true
	})
}

func (m *DevicesManager) Attach(event w.Event, listener w.Listener) {
	m.listeners.Attach(event, listener)
}

func (m *DevicesManager) DeAttach(event w.Event, listener w.Listener) {
	m.listeners.DeAttach(event, listener)
}

func (m *DevicesManager) SetTickerCheckerDuration(t time.Duration) {
	m.tickerChecker.SetDuration(t)
}

func (m *DevicesManager) Check() {
	if len(m.chanChecker) == 0 {
		m.chanChecker <- struct{}{}
	}
}

func (m *DevicesManager) doCheck() {
	for {
		select {
		case <-m.chanChecker:
			for key, device := range m.Devices() {
				// в параллель вешать нельзя, так как многие устройства висят на одной линии и
				// запросы идут последовательно, поэтому не укладывается в таймаут
				m.checker(key, device)
			}

		case <-m.tickerChecker.C():
			m.Check()
		}
	}
}

func (m *DevicesManager) checker(key string, device boggart.Device) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), m.timeoutChecker)
	defer ctxCancel()

	done := make(chan bool, 1)

	go func() {
		done <- device.Ping(ctx)
	}()

	select {
	case <-ctx.Done():
		if !device.IsEnabled() {
			return
		}

		if ctx.Err() != nil && ctx.Err() != context.Canceled {
			device.Disable()
			m.Logger().Warn("Device has been disabled because the ping failed", map[string]interface{}{
				"device.key": key,
				"device.id":  device.Id(),
				"error":      ctx.Err().Error(),
			})

			// TODO: send event
		}

		return

	case result := <-done:
		if result == device.IsEnabled() {
			return
		}

		if !result {
			device.Disable()
			m.Logger().Warn("Device has been disabled because the ping returns false", map[string]interface{}{
				"device.key": key,
				"device.id":  device.Id(),
			})

			// TODO: send event
		} else {
			device.Enable()

			m.Logger().Info("Device has been enabled because the ping returns true", map[string]interface{}{
				"device.key": key,
				"device.id":  device.Id(),
			})

			// TODO: send event
		}

		return
	}
}
