package model

import (
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"

	a "github.com/kihamo/boggart/atomic"
)

type Device struct {
	networkAddress   uint32
	ieeeAddress      *a.Bytes
	capabilities     uint32
	deviceType       uint32
	lastSeen         *a.Time
	interviewStatus  uint32
	manufacturerCode uint32
	endpoints        sync.Map
}

func NewDevice(networkAddress uint16) *Device {
	return &Device{
		networkAddress: uint32(networkAddress),
		ieeeAddress:    a.NewBytes(),
		lastSeen:       a.NewTimeDefault(time.Now()),
	}
}

func (d *Device) NetworkAddress() uint16 {
	return uint16(atomic.LoadUint32(&d.networkAddress))
}

func (d *Device) SetNetworkAddress(address uint16) {
	atomic.StoreUint32(&d.networkAddress, uint32(address))
}

func (d *Device) IEEEAddress() []byte {
	return d.ieeeAddress.Load()
}

func (d *Device) SetIEEEAddress(address []byte) {
	d.ieeeAddress.Set(address)
}

func (d *Device) IEEEAddressAsString() string {
	return hex.EncodeToString(d.IEEEAddress())
}

func (d *Device) Capabilities() uint8 {
	return uint8(atomic.LoadUint32(&d.capabilities))
}

func (d *Device) SetCapabilities(capabilities uint8) {
	atomic.StoreUint32(&d.capabilities, uint32(capabilities))
}

func (d *Device) DeviceType() uint8 {
	return uint8(atomic.LoadUint32(&d.deviceType))
}

func (d *Device) DeviceTypeAsString() string {
	switch d.DeviceType() {
	case 1:
		return "coordinator"
	case 2:
		return "router"
	case 4:
		return "end device"
	}

	return "none"
}

func (d *Device) SetDeviceType(deviceType uint8) {
	atomic.StoreUint32(&d.deviceType, uint32(deviceType))
}

func (d *Device) LastSeen() time.Time {
	return d.lastSeen.Load()
}

func (d *Device) UpdateLastSeen() {
	d.lastSeen.Set(time.Now())
}

func (d *Device) InterviewStatus() uint32 {
	return atomic.LoadUint32(&d.interviewStatus)
}

func (d *Device) SetInterviewStatus(status uint32) {
	atomic.StoreUint32(&d.interviewStatus, status)
}

func (d *Device) ManufacturerCode() uint16 {
	return uint16(atomic.LoadUint32(&d.manufacturerCode))
}

func (d *Device) SetManufacturerCode(code uint16) {
	atomic.StoreUint32(&d.manufacturerCode, uint32(code))
}

func (d *Device) EndpointAdd(endpoint *Endpoint) {
	// ignore exists
	if d.Endpoint(endpoint.ID()) != nil {
		return
	}

	d.endpoints.Store(endpoint.ID(), endpoint)
}

func (d *Device) Endpoints() []*Endpoint {
	endpoints := make([]*Endpoint, 0)

	d.endpoints.Range(func(key, value interface{}) bool {
		endpoints = append(endpoints, value.(*Endpoint))
		return true
	})

	return endpoints
}

func (d *Device) Endpoint(id uint8) *Endpoint {
	value, ok := d.endpoints.Load(id)
	if !ok {
		return nil
	}

	return value.(*Endpoint)
}
