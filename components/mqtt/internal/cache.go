package internal

import (
	"sync"

	"github.com/hashicorp/golang-lru"
	"github.com/kihamo/boggart/components/mqtt"
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

func (c *cache) Get(topic mqtt.Topic) (value []byte, ok bool) {
	c.lock.RLock()
	cached, ok := c.arc.Get(topic.String())
	c.lock.RUnlock()

	if ok {
		value = cached.([]byte)
	}

	return value, ok
}

func (c *cache) Add(topic mqtt.Topic, payload []byte) {
	c.lock.RLock()
	c.arc.Add(topic.String(), payload)
	c.lock.RUnlock()
}

func (c *cache) Payloads() map[mqtt.Topic][]byte {
	c.lock.RLock()
	cache := c.arc
	c.lock.RUnlock()

	keys := cache.Keys()
	result := make(map[mqtt.Topic][]byte, len(keys))
	for _, k := range keys {
		if v, ok := cache.Get(k); ok {
			result[mqtt.Topic(k.(string))] = v.([]byte)
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
