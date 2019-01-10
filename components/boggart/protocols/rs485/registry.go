package rs485

import (
	"sync"
	"time"
)

var (
	mutex       sync.RWMutex
	connections = make(map[string]*Connection)
)

func GetConnection(address string, timeout time.Duration) *Connection {
	mutex.RLock()
	conn, ok := connections[address]
	mutex.RUnlock()

	if !ok {
		conn = NewConnection(address, timeout)

		mutex.Lock()
		connections[address] = conn
		mutex.Unlock()
	}

	return conn
}
