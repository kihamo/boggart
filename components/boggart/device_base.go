package boggart

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"

	w "github.com/ghthor/gowol"
	"github.com/kihamo/go-workers"
	"github.com/pborman/uuid"
)

type DeviceBase struct {
	id                   atomic.Value
	description          atomic.Value
	enabled              uint64
	triggerEventsChannel chan DeviceTriggerEvent
}

func (d *DeviceBase) Init() {
	d.triggerEventsChannel = make(chan DeviceTriggerEvent)
	d.SetId(uuid.New())
	d.Enable()
}

func (d *DeviceBase) Id() string {
	var id string

	if value := d.id.Load(); value != nil {
		id = value.(string)
	}

	return id
}

func (d *DeviceBase) SetId(id string) {
	d.id.Store(id)
}

func (d *DeviceBase) Description() string {
	var description string

	if value := d.description.Load(); value != nil {
		description = value.(string)
	}

	return description
}

func (d *DeviceBase) SetDescription(description string) {
	d.description.Store(description)
}

func (d *DeviceBase) Types() []DeviceType {
	return nil
}

func (d *DeviceBase) IsEnabled() bool {
	return atomic.LoadUint64(&d.enabled) == 1
}

func (d *DeviceBase) Enable() error {
	atomic.StoreUint64(&d.enabled, 1)
	d.TriggerEvent(DeviceEventDeviceEnabled, d)

	return nil
}

func (d *DeviceBase) Ping(_ context.Context) bool {
	return false
}

func (d *DeviceBase) Disable() error {
	atomic.StoreUint64(&d.enabled, 0)
	d.TriggerEvent(DeviceEventDeviceDisabled, d)

	return nil
}

func (d *DeviceBase) Listeners() []workers.ListenerWithEvents {
	return nil
}

func (d *DeviceBase) Tasks() []workers.Task {
	return nil
}

func (d *DeviceBase) TriggerEventChannel() <-chan DeviceTriggerEvent {
	return d.triggerEventsChannel
}

func (d *DeviceBase) TriggerEvent(event workers.Event, args ...interface{}) {
	if d.triggerEventsChannel == nil {
		return
	}

	go func() {
		d.triggerEventsChannel <- NewDeviceTriggerEventBase(event, append([]interface{}{d}, args...))
	}()
}

type DeviceSerialNumber struct {
	mutex        sync.RWMutex
	serialNumber string
}

func (d *DeviceSerialNumber) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *DeviceSerialNumber) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type DeviceWOL struct {
	mutex  sync.RWMutex
	mac    net.HardwareAddr
	ip     net.IP
	subnet net.IP
}

func (d *DeviceWOL) Mac() net.HardwareAddr {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.mac
}

func (d *DeviceWOL) SetMac(mac string) error {
	parsed, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}

	d.mutex.Lock()
	d.mac = parsed
	d.mutex.Unlock()

	return nil
}

func (d *DeviceWOL) IP() net.IP {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.ip
}

func (d *DeviceWOL) SetIP(ip string) {
	d.mutex.Lock()
	d.ip = net.ParseIP(ip)
	d.mutex.Unlock()
}

func (d *DeviceWOL) Subnet() net.IP {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.subnet
}

func (d *DeviceWOL) SetSubnet(subnet string) {
	d.mutex.Lock()
	d.subnet = net.ParseIP(subnet)
	d.mutex.Unlock()
}

func (d *DeviceWOL) WakeUp() error {
	mac := d.Mac()
	if mac == nil {
		return errors.New("Mac isn't set")
	}

	var broadcastAddress net.IP

	ip := d.IP()
	subnet := d.Subnet()
	if ip != nil && subnet != nil {
		broadcastAddress = net.IP{0, 0, 0, 0}
		for i := 0; i < 4; i++ {
			broadcastAddress[i] = (ip[i] & subnet[i]) | ^subnet[i]
		}
	} else {
		broadcastAddress = net.IP{255, 255, 255, 255}
	}

	return w.MagicWake(mac.String(), broadcastAddress.String())
}
