package internal

import (
	"sync"
	"time"

	"github.com/hashicorp/golang-lru"
	"github.com/kihamo/boggart/components/mqtt"
)

type cacheItem struct {
	datetime time.Time
	topic    mqtt.Topic
	payload  []byte
}

type cache struct {
	arc  *lru.ARCCache
	lock sync.RWMutex
}

func newCacheItem(topic mqtt.Topic, payload []byte) *cacheItem {
	return &cacheItem{
		datetime: time.Now(),
		topic:    topic,
		payload:  payload,
	}
}

func (i *cacheItem) Payload() []byte {
	return i.payload
}

func (i *cacheItem) Datetime() time.Time {
	return i.datetime
}

func (i *cacheItem) Topic() mqtt.Topic {
	return i.topic
}

func newCache(size int) (*cache, error) {
	c := &cache{}

	if err := c.Resize(size); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *cache) Get(topic mqtt.Topic) (item mqtt.CacheItem, ok bool) {
	c.lock.RLock()
	cached, ok := c.arc.Get(topic.String())
	c.lock.RUnlock()

	if ok {
		item = cached.(*cacheItem)
	}

	return item, ok
}

func (c *cache) Add(topic mqtt.Topic, payload []byte) {
	c.lock.RLock()
	c.arc.Add(topic.String(), newCacheItem(topic, payload))
	c.lock.RUnlock()
}

func (c *cache) Payloads() []mqtt.CacheItem {
	c.lock.RLock()
	cache := c.arc
	c.lock.RUnlock()

	keys := cache.Keys()
	result := make([]mqtt.CacheItem, 0, len(keys))
	for _, k := range keys {
		if v, ok := cache.Get(k); ok {
			result = append(result, v.(*cacheItem))
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
