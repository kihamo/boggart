package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

const (
	DefaultTimeoutChecker = time.Second * 2

	devicesManagerNotReady = int64(iota)
	devicesManagerReady
)

type DevicesManager struct {
	mutex sync.RWMutex

	ready          int64
	storage        *sync.Map
	mqtt           mqtt.Component
	workers        workers.Component
	listeners      *manager.ListenersManager
	chanChecker    chan string
	tickerChecker  *w.Ticker
	timeoutChecker time.Duration
}

func NewDevicesManager(mqtt mqtt.Component, workers workers.Component, listeners *manager.ListenersManager) *DevicesManager {
	m := &DevicesManager{
		ready:          devicesManagerNotReady,
		storage:        new(sync.Map),
		mqtt:           mqtt,
		workers:        workers,
		listeners:      listeners,
		chanChecker:    make(chan string),
		tickerChecker:  w.NewTicker(time.Minute),
		timeoutChecker: DefaultTimeoutChecker,
	}

	go m.doCheck()

	return m
}

func (m *DevicesManager) Register(device boggart.Device) string {
	id := uuid.New()
	m.RegisterWithID(id, device)
	return id
}

func (m *DevicesManager) RegisterWithID(id string, device boggart.Device) {
	m.storage.Store(id, device)
	m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDeviceRegister, device, id)

	if subs, ok := device.(boggart.DeviceHasMQTTSubscribers); ok {
		m.mqtt.SubscribeSubscribers(subs.MQTTSubscribers())
	}

	if tasks, ok := device.(boggart.DeviceHasTasks); ok {
		for _, task := range tasks.Tasks() {
			m.workers.AddTask(task)
		}
	}

	if listeners, ok := device.(boggart.DeviceHasListeners); ok {
		for _, listener := range listeners.Listeners() {
			m.listeners.AddListener(listener)
		}
	}

	events := device.TriggerEventChannel()
	if events != nil {
		go m.doDeviceEvents(events)
	}

	m.CheckByKeys(id)
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

func (m *DevicesManager) Ready() {
	if !m.IsReady() {
		atomic.StoreInt64(&m.ready, devicesManagerReady)
		m.listeners.AsyncTrigger(context.TODO(), boggart.DeviceEventDevicesManagerReady)

		// TODO: запускать рутину автоматического опроса только после завершения инициализации
	}
}

func (m *DevicesManager) IsReady() bool {
	return atomic.LoadInt64(&m.ready) == devicesManagerReady
}

func (m *DevicesManager) SetCheckerTickerDuration(duration time.Duration) {
	m.tickerChecker.SetDuration(duration)
}

func (m *DevicesManager) SetCheckerTimeout(duration time.Duration) {
	m.mutex.Lock()
	m.timeoutChecker = duration
	m.mutex.Unlock()
}

func (m *DevicesManager) Check() {
	keys := make([]string, 0, 0)

	m.storage.Range(func(key interface{}, device interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})

	m.CheckByKeys(keys...)
}

func (m *DevicesManager) CheckByKeys(keys ...string) {
	go func() {
		for _, key := range keys {
			m.chanChecker <- key
		}
	}()
}

func (m *DevicesManager) doCheck() {
	for {
		select {
		case key := <-m.chanChecker:
			if device := m.Device(key); device != nil {
				m.checker(key, device)
			}

		case <-m.tickerChecker.C():
			m.Check()
		}
	}
}

func (m *DevicesManager) doDeviceEvents(ch <-chan boggart.DeviceTriggerEvent) {
	for event := range ch {
		m.listeners.AsyncTrigger(event.Context(), event.Event(), event.Arguments()...)
	}
}

func (m *DevicesManager) checker(key string, device boggart.Device) {
	m.mutex.RLock()
	timeout := m.timeoutChecker
	m.mutex.RUnlock()

	ctx, ctxCancel := context.WithTimeout(context.TODO(), timeout)
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
			if err := device.Disable(); err == nil {
				m.listeners.AsyncTrigger(ctx, boggart.DeviceEventDeviceDisabledAfterCheck, device, key, ctx.Err())
			}
		}

		return

	case result := <-done:
		if result == device.IsEnabled() {
			return
		}

		if !result {
			if err := device.Disable(); err == nil {
				m.listeners.AsyncTrigger(ctx, boggart.DeviceEventDeviceDisabledAfterCheck, device, key, nil)
			}
		} else {
			if err := device.Enable(); err == nil {
				m.listeners.AsyncTrigger(ctx, boggart.DeviceEventDeviceEnabledAfterCheck, device, key)
			}
		}

		return
	}
}
