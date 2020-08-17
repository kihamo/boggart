package model

import (
	"sync"
	"sync/atomic"
)

type Endpoint struct {
	id             uint32
	profileID      uint32
	inClusterList  []uint16
	outClusterList []uint16

	mutex sync.RWMutex
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

func (e *Endpoint) InClusterList() []uint16 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return append([]uint16(nil), e.inClusterList...)
}

func (e *Endpoint) SetInClusterList(list []uint16) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.inClusterList = list
}

func (e *Endpoint) OutClusterList() []uint16 {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return append([]uint16(nil), e.outClusterList...)
}

func (e *Endpoint) SetOutClusterList(list []uint16) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.outClusterList = list
}
