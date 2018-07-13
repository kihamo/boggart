package internal

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
	"github.com/kihamo/go-workers/manager"
)

type SecurityManager struct {
	listener.BaseListener

	// примитивы должны быть первыми при объявлении для arm систем
	status         int64
	devicesManager boggart.DevicesManager
	listeners      *manager.ListenersManager
}

func NewSecurityManager(devicesManager boggart.DevicesManager, listeners *manager.ListenersManager) *SecurityManager {
	s := &SecurityManager{
		devicesManager: devicesManager,
		listeners:      listeners,
		status:         int64(boggart.SecurityStatusOpen),
	}
	s.Init()
	s.SetName(boggart.ComponentName + ".security")
	listeners.AddListener(s)

	return s
}

func (s *SecurityManager) Status() boggart.SecurityStatus {
	return boggart.SecurityStatus(atomic.LoadInt64(&s.status))
}

func (s *SecurityManager) setStatus(status boggart.SecurityStatus) {
	prev := s.IsOpen()
	atomic.StoreInt64(&s.status, int64(status))
	current := s.IsOpen()

	if prev != current {
		if current {
			s.listeners.AsyncTrigger(boggart.SecurityOpen, status)
		} else {
			s.listeners.AsyncTrigger(boggart.SecurityClosed, status)
		}
	}
}

func (s *SecurityManager) IsClosed() bool {
	return s.Status() == boggart.SecurityStatusClosedForce || s.Status() == boggart.SecurityStatusClosed
}

func (s *SecurityManager) IsOpen() bool {
	return !s.IsClosed()
}

func (s *SecurityManager) IsForce() bool {
	return s.Status() == boggart.SecurityStatusOpenForce || s.Status() == boggart.SecurityStatusClosedForce
}

func (s *SecurityManager) Close() {
	s.setStatus(boggart.SecurityStatusClosed)
}

func (s *SecurityManager) Open() {
	s.setStatus(boggart.SecurityStatusOpen)
}

func (s *SecurityManager) CloseForce() {
	s.setStatus(boggart.SecurityStatusClosedForce)
}

func (s *SecurityManager) OpenForce() {
	s.setStatus(boggart.SecurityStatusOpenForce)
}

func (s *SecurityManager) Events() []workers.Event {
	return []workers.Event{
		boggart.SecurityClosed,
		boggart.DeviceEventDevicesManagerReady,
		boggart.DeviceEventWifiClientConnected,
		boggart.DeviceEventWifiClientDisconnected,
		boggart.DeviceEventHikvisionEventNotificationAlert,
		boggart.DeviceEventDeviceEnabled,
		boggart.DeviceEventDeviceEnabledAfterCheck,
		devices.EventDoorGPIOReedSwitchOpen,
		devices.EventDoorGPIOReedSwitchClose,
	}
}

func (s *SecurityManager) Run(_ context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case devices.EventDoorGPIOReedSwitchOpen:
		if s.IsOpen() {
			if s.checkClosed() {
				s.Close()
			}
		} else {
			s.Open()
		}

	case devices.EventDoorGPIOReedSwitchClose:
		if s.IsOpen() {
			if s.checkClosed() {
				s.Close()
			}
		} else {
			s.Open()
		}

	case boggart.SecurityClosed:
		// выключаем устройства, которые не должны работать в закрытом контуре
		activeDevices := s.devicesManager.DevicesByTypes([]boggart.DeviceType{
			boggart.DeviceTypeTV,
			boggart.DeviceTypeLight,
		})
		for _, device := range activeDevices {
			if device.IsEnabled() {
				device.Disable()
			}
		}

	case boggart.SecurityOpen:
		// TODO: снимаем с охраны камеры

	// после опроса устройств устанавливаем текущий статус контура
	case boggart.DeviceEventDevicesManagerReady:
		switch s.Status() {
		case boggart.SecurityStatusClosedForce, boggart.SecurityStatusOpenForce:
			return
		}

		if s.checkClosed() {
			s.Close()
		} else {
			s.Open()
		}
	}
}

// Подготовить объект к приходу через заданное время
func (s *SecurityManager) PrepareOpen(interval time.Duration) {
	if s.IsOpen() {
		return
	}

	s.listeners.AsyncTrigger(boggart.SecurityPrepareOpen, interval)
}

// Проверяет контур на закрытость
func (s *SecurityManager) checkClosed() bool {
	// для ручного управления проверки игнорируются
	switch s.Status() {
	case boggart.SecurityStatusOpenForce:
		return false
	case boggart.SecurityStatusClosedForce:
		return true
	}

	// если подключены устройства к вайфаю, то в квартире кто-то есть
	// если зафиксирован расход воды, то в квартире кто-то есть (расход должен быть значительный)
	// если зафиксирован значительный расход энергии то в квартире кто-то есть

	// если хотя бы один телевизор включен, то в квартире кто-то есть
	activeDevices := s.devicesManager.DevicesByTypes([]boggart.DeviceType{boggart.DeviceTypeTV})
	for _, device := range activeDevices {
		if device.IsEnabled() {
			return false
		}
	}

	return true
}
