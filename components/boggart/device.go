package boggart

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/go-workers"
	"github.com/kihamo/snitch"
	"github.com/pborman/uuid"
)

type DeviceType int64

const (
	DeviceTypeCamera DeviceType = iota
	DeviceTypePhone
	DeviceTypeVideoRecorder
)

type DeviceManager interface {
	snitch.Collector

	Register(string, Device)
	Device(string) Device
	Devices() map[string]Device
	DevicesByTypes([]DeviceType) map[string]Device
}

type Device interface {
	Id() string
	Description() string
	Types() []DeviceType
	IsEnabled() bool
	Disable()
	Enable()
	Tasks() []workers.Task
}

type ReedSwitch interface {
	Device

	IsOpen() bool
	IsClose() bool
}

type Camera interface {
	Device

	Snapshot(context.Context) ([]byte, error)
}

type Phone interface {
	Device

	Number() string
	Balance(context.Context) (float64, error)
}

type ElectricityMeter interface {
	Device
}

type WaterMeter interface {
	Device
}

type HeatMeter interface {
	Device
}

type Router interface {
	Device
}

type Thermometer interface {
	Device
}

type Barometer interface {
	Device
}

type WaterDetector interface {
	Device
}

type MotionDetector interface {
	Device
}

type DeviceBase struct {
	id          atomic.Value
	description atomic.Value
	enabled     uint64
}

func (d *DeviceBase) Init() {
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

func (d *DeviceBase) Enable() {
	atomic.StoreUint64(&d.enabled, 1)
}

func (d *DeviceBase) Disable() {
	atomic.StoreUint64(&d.enabled, 0)
}

func (d *DeviceBase) Tasks() []workers.Task {
	return nil
}
