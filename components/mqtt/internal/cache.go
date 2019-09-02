package internal

import (
	"sync"

	"github.com/hashicorp/golang-lru"
)

type cache struct {
	arc  *lru.ARCCache
	lock sync.RWMutex
}

func newCache(size int) (*cache, error) {
	c := &cache{}

	if err := c.Resize(size); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *cache) Get(topic string) (value []byte, ok bool) {
	c.lock.RLock()
	cached, ok := c.arc.Get(topic)
	c.lock.RUnlock()

	if ok {
		value = cached.([]byte)
	}

	return value, ok
}

func (c *cache) Add(topic string, payload []byte) {
	c.lock.RLock()
	c.arc.Add(topic, payload)
	c.lock.RUnlock()
}

func (c *cache) Payloads() map[string][]byte {
	c.lock.RLock()
	cache := c.arc
	c.lock.RUnlock()

	keys := cache.Keys()
	result := make(map[string][]byte, len(keys))
	for _, k := range keys {
		if v, ok := cache.Get(k); ok {
			result[k.(string)] = v.([]byte)
		}
	}

	return result
}

func (c *cache) Resize(size int) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	arc, err := lru.NewARC(size)
	if err != nil {
		return err
	}

	c.arc = arc
	return nil
}

func (c *cache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.arc.Len()
}

func (c *cache) Purge() {
	c.lock.RLock()
	c.arc.Purge()
	c.lock.RUnlock()
}
