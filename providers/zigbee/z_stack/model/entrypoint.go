package model

import (
	"sync/atomic"
)

type Endpoint struct {
	id             uint32
	profileID      uint32
	deviceID       uint16
	inClusterList  []uint16
	outClusterList []uint16
}

func NewEndpoint(id uint8) *Endpoint {
	return &Endpoint{
		id: uint32(id),
	}
}

func (e *Endpoint) ID() uint8 {
	return uint8(atomic.LoadUint32(&e.id))
}

func (e *Endpoint) ProfileID() uint16 {
	return uint16(atomic.LoadUint32(&e.profileID))
}

func (e *Endpoint) SetProfileID(id uint16) {
	atomic.StoreUint32(&e.profileID, uint32(id))
}
